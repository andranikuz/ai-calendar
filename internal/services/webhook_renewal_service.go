package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/andranikuz/smart-goal-calendar/internal/adapters/google"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
)

// WebhookRenewalService handles automatic renewal of Google Calendar webhook subscriptions
type WebhookRenewalService struct {
	calendarService        *google.CalendarService
	oauth2Service          *google.OAuth2Service
	googleIntegrationRepo  repositories.GoogleIntegrationRepository
	googleCalendarSyncRepo repositories.GoogleCalendarSyncRepository
	
	// Configuration
	renewalCheckInterval time.Duration // How often to check for expiring webhooks
	renewalThreshold     time.Duration // How early to renew before expiration
	
	// Control
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.RWMutex
}

// NewWebhookRenewalService creates a new webhook renewal service
func NewWebhookRenewalService(
	calendarService *google.CalendarService,
	oauth2Service *google.OAuth2Service,
	googleIntegrationRepo repositories.GoogleIntegrationRepository,
	googleCalendarSyncRepo repositories.GoogleCalendarSyncRepository,
) *WebhookRenewalService {
	return &WebhookRenewalService{
		calendarService:        calendarService,
		oauth2Service:          oauth2Service,
		googleIntegrationRepo:  googleIntegrationRepo,
		googleCalendarSyncRepo: googleCalendarSyncRepo,
		
		// Check every hour
		renewalCheckInterval: time.Hour,
		// Renew 24 hours before expiration
		renewalThreshold:     24 * time.Hour,
		
		stopChan: make(chan struct{}),
	}
}

// Start begins the webhook renewal background service
func (s *WebhookRenewalService) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.running {
		log.Println("WebhookRenewalService: Already running")
		return
	}
	
	s.running = true
	s.wg.Add(1)
	
	go s.renewalLoop(ctx)
	log.Println("WebhookRenewalService: Started")
}

// Stop gracefully stops the webhook renewal service
func (s *WebhookRenewalService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if !s.running {
		return
	}
	
	close(s.stopChan)
	s.wg.Wait()
	s.running = false
	
	log.Println("WebhookRenewalService: Stopped")
}

// IsRunning returns whether the service is currently running
func (s *WebhookRenewalService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// renewalLoop is the main loop that checks and renews expiring webhooks
func (s *WebhookRenewalService) renewalLoop(ctx context.Context) {
	defer s.wg.Done()
	
	ticker := time.NewTicker(s.renewalCheckInterval)
	defer ticker.Stop()
	
	// Do initial check immediately
	s.checkAndRenewWebhooks(ctx)
	
	for {
		select {
		case <-ctx.Done():
			log.Println("WebhookRenewalService: Context cancelled")
			return
		case <-s.stopChan:
			log.Println("WebhookRenewalService: Stop signal received")
			return
		case <-ticker.C:
			s.checkAndRenewWebhooks(ctx)
		}
	}
}

// checkAndRenewWebhooks finds and renews expiring webhook subscriptions
func (s *WebhookRenewalService) checkAndRenewWebhooks(ctx context.Context) {
	log.Println("WebhookRenewalService: Checking for expiring webhooks")
	
	// Get all calendar sync configurations with active webhooks
	syncs, err := s.googleCalendarSyncRepo.GetActiveWebhooks(ctx)
	if err != nil {
		log.Printf("WebhookRenewalService: Error getting active webhooks: %v", err)
		return
	}
	
	log.Printf("WebhookRenewalService: Found %d active webhook subscriptions", len(syncs))
	
	now := time.Now()
	renewedCount := 0
	errorCount := 0
	
	for _, sync := range syncs {
		// Check if webhook needs renewal
		if s.needsRenewal(sync, now) {
			log.Printf("WebhookRenewalService: Renewing webhook for calendar %s (expires: %v)", 
				sync.CalendarID, sync.WebhookExpiresAt)
			
			err := s.renewWebhook(ctx, sync)
			if err != nil {
				log.Printf("WebhookRenewalService: Failed to renew webhook for calendar %s: %v", 
					sync.CalendarID, err)
				errorCount++
				
				// Update sync with error info
				s.updateSyncError(ctx, sync, fmt.Sprintf("Webhook renewal failed: %v", err))
			} else {
				log.Printf("WebhookRenewalService: Successfully renewed webhook for calendar %s", 
					sync.CalendarID)
				renewedCount++
			}
		}
	}
	
	if renewedCount > 0 || errorCount > 0 {
		log.Printf("WebhookRenewalService: Renewal complete. Renewed: %d, Errors: %d", 
			renewedCount, errorCount)
	}
}

// needsRenewal checks if a webhook subscription needs to be renewed
func (s *WebhookRenewalService) needsRenewal(sync *entities.GoogleCalendarSync, now time.Time) bool {
	// If no expiry time is set, we can't determine if it needs renewal
	if sync.WebhookExpiresAt == nil {
		return false
	}
	
	// If webhook channel ID is empty, it's not an active webhook
	if sync.WebhookChannelID == "" {
		return false
	}
	
	// Check if webhook expires within the renewal threshold
	timeUntilExpiry := sync.WebhookExpiresAt.Sub(now)
	return timeUntilExpiry <= s.renewalThreshold
}

// renewWebhook renews a specific webhook subscription
func (s *WebhookRenewalService) renewWebhook(ctx context.Context, sync *entities.GoogleCalendarSync) error {
	// Get the Google integration for authentication
	integration, err := s.googleIntegrationRepo.GetByID(ctx, sync.GoogleIntegrationID)
	if err != nil {
		return fmt.Errorf("failed to get Google integration: %w", err)
	}
	
	if integration == nil || !integration.Enabled {
		return fmt.Errorf("Google integration not found or disabled")
	}
	
	// Refresh access token if needed
	accessToken := integration.AccessToken
	if integration.IsTokenExpiringSoon() {
		newTokens, err := s.oauth2Service.RefreshToken(ctx, integration.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		
		// Update tokens in database
		err = s.googleIntegrationRepo.UpdateTokens(ctx, integration.ID, 
			newTokens.AccessToken, newTokens.RefreshToken, newTokens.Expiry)
		if err != nil {
			return fmt.Errorf("failed to update tokens: %w", err)
		}
		
		accessToken = newTokens.AccessToken
	}
	
	// Stop the existing webhook (if it exists and is still active)
	if sync.WebhookChannelID != "" && sync.WebhookResourceID != "" {
		err = s.calendarService.StopWebhook(ctx, accessToken, sync.WebhookChannelID, sync.WebhookResourceID)
		if err != nil {
			// Log error but don't fail the renewal - the old webhook might already be expired
			log.Printf("WebhookRenewalService: Warning: Failed to stop old webhook %s: %v", 
				sync.WebhookChannelID, err)
		}
	}
	
	// Create a new webhook subscription
	newWebhook, err := s.calendarService.SetupWebhookWithExpiry(ctx, accessToken, sync.CalendarID, sync.WebhookURL)
	if err != nil {
		return fmt.Errorf("failed to setup new webhook: %w", err)
	}
	
	// Update the sync configuration with new webhook details
	now := time.Now()
	sync.WebhookChannelID = newWebhook.ChannelID
	sync.WebhookResourceID = newWebhook.ResourceID
	sync.WebhookExpiresAt = &newWebhook.ExpiresAt
	sync.LastSyncError = "" // Clear any previous errors
	sync.UpdatedAt = now
	
	err = s.googleCalendarSyncRepo.Update(ctx, sync)
	if err != nil {
		return fmt.Errorf("failed to update sync configuration: %w", err)
	}
	
	return nil
}

// updateSyncError updates the sync configuration with error information
func (s *WebhookRenewalService) updateSyncError(ctx context.Context, sync *entities.GoogleCalendarSync, errorMsg string) {
	sync.LastSyncError = errorMsg
	sync.UpdatedAt = time.Now()
	
	err := s.googleCalendarSyncRepo.Update(ctx, sync)
	if err != nil {
		log.Printf("WebhookRenewalService: Failed to update sync error: %v", err)
	}
}

// ForceRenewWebhook manually renews a specific webhook (for testing or manual intervention)
func (s *WebhookRenewalService) ForceRenewWebhook(ctx context.Context, syncID string) error {
	sync, err := s.googleCalendarSyncRepo.GetByID(ctx, syncID)
	if err != nil {
		return fmt.Errorf("failed to get sync configuration: %w", err)
	}
	
	if sync == nil {
		return fmt.Errorf("sync configuration not found")
	}
	
	return s.renewWebhook(ctx, sync)
}

// GetExpiringWebhooks returns webhooks that will expire within the specified duration
func (s *WebhookRenewalService) GetExpiringWebhooks(ctx context.Context, within time.Duration) ([]*entities.GoogleCalendarSync, error) {
	syncs, err := s.googleCalendarSyncRepo.GetActiveWebhooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active webhooks: %w", err)
	}
	
	var expiring []*entities.GoogleCalendarSync
	now := time.Now()
	
	for _, sync := range syncs {
		if sync.WebhookExpiresAt != nil {
			timeUntilExpiry := sync.WebhookExpiresAt.Sub(now)
			if timeUntilExpiry <= within && timeUntilExpiry > 0 {
				expiring = append(expiring, sync)
			}
		}
	}
	
	return expiring, nil
}

// WebhookInfo contains information about a webhook subscription
type WebhookInfo struct {
	ChannelID   string    `json:"channel_id"`
	ResourceID  string    `json:"resource_id"`
	ExpiresAt   time.Time `json:"expires_at"`
}