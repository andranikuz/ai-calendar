package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/andranikuz/smart-goal-calendar/internal/adapters/google"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
)

type GoogleWebhookHandler struct {
	calendarService        *google.CalendarService
	oauth2Service          *google.OAuth2Service
	googleIntegrationRepo  repositories.GoogleIntegrationRepository
	googleCalendarSyncRepo repositories.GoogleCalendarSyncRepository
	eventRepo              repositories.EventRepository
}

func NewGoogleWebhookHandler(
	calendarService *google.CalendarService,
	oauth2Service *google.OAuth2Service,
	googleIntegrationRepo repositories.GoogleIntegrationRepository,
	googleCalendarSyncRepo repositories.GoogleCalendarSyncRepository,
	eventRepo repositories.EventRepository,
) *GoogleWebhookHandler {
	return &GoogleWebhookHandler{
		calendarService:        calendarService,
		oauth2Service:          oauth2Service,
		googleIntegrationRepo:  googleIntegrationRepo,
		googleCalendarSyncRepo: googleCalendarSyncRepo,
		eventRepo:              eventRepo,
	}
}

// GoogleWebhookNotification represents the webhook notification from Google Calendar
type GoogleWebhookNotification struct {
	ChannelID     string `json:"channelId"`
	ChannelToken  string `json:"channelToken,omitempty"`
	ChannelExpiry string `json:"channelExpiry,omitempty"`
	ResourceID    string `json:"resourceId"`
	ResourceURI   string `json:"resourceUri"`
	ResourceState string `json:"resourceState"`
	MessageNumber string `json:"messageNumber,omitempty"`
}

// HandleCalendarWebhook processes Google Calendar webhook notifications
func (h *GoogleWebhookHandler) HandleCalendarWebhook(c *gin.Context) {
	// Extract Google-specific headers
	channelID := c.GetHeader("X-Goog-Channel-Id")
	channelToken := c.GetHeader("X-Goog-Channel-Token")
	channelExpiry := c.GetHeader("X-Goog-Channel-Expiry")
	resourceID := c.GetHeader("X-Goog-Resource-Id")
	resourceURI := c.GetHeader("X-Goog-Resource-Uri")
	resourceState := c.GetHeader("X-Goog-Resource-State")
	messageNumber := c.GetHeader("X-Goog-Message-Number")

	// Log webhook details for debugging
	fmt.Printf("Webhook received: channelID=%s, resourceState=%s, resourceURI=%s\n", 
		channelID, resourceState, resourceURI)

	// Handle sync state (Google's webhook verification)
	if resourceState == "sync" {
		c.Status(http.StatusOK)
		return
	}

	// Validate required headers
	if channelID == "" || resourceID == "" || resourceURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required webhook headers"})
		return
	}

	notification := &GoogleWebhookNotification{
		ChannelID:     channelID,
		ChannelToken:  channelToken,
		ChannelExpiry: channelExpiry,
		ResourceID:    resourceID,
		ResourceURI:   resourceURI,
		ResourceState: resourceState,
		MessageNumber: messageNumber,
	}

	// Process the webhook notification asynchronously
	go h.processWebhookNotification(context.Background(), notification)

	c.Status(http.StatusOK)
}

// processWebhookNotification handles the actual webhook processing
func (h *GoogleWebhookHandler) processWebhookNotification(ctx context.Context, notification *GoogleWebhookNotification) {
	// Extract calendar ID from resource URI
	calendarID := h.extractCalendarIDFromURI(notification.ResourceURI)
	if calendarID == "" {
		fmt.Printf("Could not extract calendar ID from URI: %s\n", notification.ResourceURI)
		return
	}

	// Find the calendar sync configuration for this channel
	sync, err := h.googleCalendarSyncRepo.GetByChannelID(ctx, notification.ChannelID)
	if err != nil {
		fmt.Printf("Error finding sync for channel %s: %v\n", notification.ChannelID, err)
		return
	}

	if sync == nil {
		fmt.Printf("No sync configuration found for channel: %s\n", notification.ChannelID)
		return
	}

	// Get the Google integration for authentication
	integration, err := h.googleIntegrationRepo.GetByID(ctx, sync.GoogleIntegrationID)
	if err != nil {
		fmt.Printf("Error getting Google integration: %v\n", err)
		return
	}

	if integration == nil || !integration.Enabled {
		fmt.Printf("Google integration not found or disabled for sync: %s\n", sync.ID)
		return
	}

	// Refresh token if needed
	accessToken := integration.AccessToken
	if integration.IsTokenExpiringSoon() {
		newTokens, err := h.oauth2Service.RefreshToken(ctx, integration.RefreshToken)
		if err != nil {
			fmt.Printf("Error refreshing token: %v\n", err)
			return
		}

		// Update tokens in database
		err = h.googleIntegrationRepo.UpdateTokens(ctx, integration.ID, 
			newTokens.AccessToken, newTokens.RefreshToken, newTokens.Expiry)
		if err != nil {
			fmt.Printf("Error updating tokens: %v\n", err)
			return
		}

		accessToken = newTokens.AccessToken
	}

	// Perform incremental sync based on the notification
	err = h.performIncrementalSync(ctx, sync, accessToken, calendarID)
	if err != nil {
		fmt.Printf("Error performing incremental sync: %v\n", err)
		return
	}

	// Update last sync time
	sync.LastSyncAt = time.Now()
	err = h.googleCalendarSyncRepo.Update(ctx, sync)
	if err != nil {
		fmt.Printf("Error updating sync timestamp: %v\n", err)
	}

	fmt.Printf("Successfully processed webhook for calendar: %s\n", calendarID)
}

// performIncrementalSync syncs changes from Google Calendar
func (h *GoogleWebhookHandler) performIncrementalSync(ctx context.Context, sync *entities.GoogleCalendarSync, accessToken, calendarID string) error {
	// Get events from Google Calendar since last sync
	var timeMin time.Time
	if sync.LastSyncAt.IsZero() {
		// If never synced, get events from last 30 days
		timeMin = time.Now().AddDate(0, 0, -30)
	} else {
		// Get events since last sync
		timeMin = sync.LastSyncAt
	}

	events, err := h.calendarService.GetEvents(ctx, accessToken, calendarID, timeMin, time.Now().AddDate(0, 3, 0))
	if err != nil {
		return fmt.Errorf("failed to get events from Google Calendar: %w", err)
	}

	// Process each event - convert to CalendarEvent format for processing
	for _, googleEvent := range events {
		// Convert GoogleCalendarEvent to CalendarEvent for processing
		calendarEvent := h.convertToCalendarEvent(googleEvent)
		err := h.syncGoogleEvent(ctx, sync, calendarEvent)
		if err != nil {
			fmt.Printf("Error syncing event %s: %v\n", googleEvent.ID, err)
			// Continue with other events
		}
	}

	return nil
}

// convertToCalendarEvent converts GoogleCalendarEvent to CalendarEvent
func (h *GoogleWebhookHandler) convertToCalendarEvent(googleEvent *google.GoogleCalendarEvent) *google.CalendarEvent {
	calendarEvent := &google.CalendarEvent{
		ID:          googleEvent.ID,
		Summary:     googleEvent.Summary,
		Description: googleEvent.Description,
		Location:    googleEvent.Location,
		Status:      googleEvent.Status,
		Attendees:   googleEvent.Attendees,
	}

	// Convert time fields
	calendarEvent.Start.DateTime = googleEvent.StartTime.Format(time.RFC3339)
	calendarEvent.End.DateTime = googleEvent.EndTime.Format(time.RFC3339)
	
	if googleEvent.AllDay {
		calendarEvent.Start.Date = googleEvent.StartTime.Format("2006-01-02")
		calendarEvent.End.Date = googleEvent.EndTime.Format("2006-01-02")
		calendarEvent.Start.DateTime = ""
		calendarEvent.End.DateTime = ""
	}

	return calendarEvent
}

// syncGoogleEvent synchronizes a single Google Calendar event
func (h *GoogleWebhookHandler) syncGoogleEvent(ctx context.Context, sync *entities.GoogleCalendarSync, googleEvent *google.CalendarEvent) error {
	// Check if event already exists in our database
	existingEvent, err := h.eventRepo.GetByGoogleEventID(ctx, googleEvent.ID)
	if err != nil && err.Error() != "event not found" {
		return fmt.Errorf("error checking existing event: %w", err)
	}

	// Handle deleted events
	if googleEvent.Status == "cancelled" {
		if existingEvent != nil {
			return h.eventRepo.Delete(ctx, existingEvent.ID)
		}
		return nil // Event doesn't exist in our system, nothing to delete
	}

	// Convert Google event to our event entity
	event, err := h.convertGoogleEventToEvent(sync.UserID, googleEvent)
	if err != nil {
		return fmt.Errorf("error converting Google event: %w", err)
	}

	if existingEvent != nil {
		// Update existing event
		event.ID = existingEvent.ID
		event.CreatedAt = existingEvent.CreatedAt
		event.UpdatedAt = time.Now()
		return h.eventRepo.Update(ctx, event)
	} else {
		// Create new event
		event.ID = entities.EventID(uuid.New().String())
		now := time.Now()
		event.CreatedAt = now
		event.UpdatedAt = now
		return h.eventRepo.Create(ctx, event)
	}
}

// convertGoogleEventToEvent converts a Google Calendar event to our internal event format
func (h *GoogleWebhookHandler) convertGoogleEventToEvent(userID entities.UserID, googleEvent *google.CalendarEvent) (*entities.Event, error) {
	startTime, err := time.Parse(time.RFC3339, googleEvent.Start.DateTime)
	if err != nil {
		// Handle all-day events
		if googleEvent.Start.Date != "" {
			startTime, err = time.Parse("2006-01-02", googleEvent.Start.Date)
			if err != nil {
				return nil, fmt.Errorf("invalid start date: %w", err)
			}
		} else {
			return nil, fmt.Errorf("invalid start time: %w", err)
		}
	}

	endTime, err := time.Parse(time.RFC3339, googleEvent.End.DateTime)
	if err != nil {
		// Handle all-day events
		if googleEvent.End.Date != "" {
			endTime, err = time.Parse("2006-01-02", googleEvent.End.Date)
			if err != nil {
				return nil, fmt.Errorf("invalid end date: %w", err)
			}
		} else {
			return nil, fmt.Errorf("invalid end time: %w", err)
		}
	}

	// Convert attendees
	var attendees []entities.Attendee
	for _, attendee := range googleEvent.Attendees {
		attendees = append(attendees, entities.Attendee{
			Email:  attendee.Email,
			Name:   attendee.Email, // Use email as name if no display name provided
			Status: entities.AttendeeStatusPending,
		})
	}

	// Determine timezone
	timezone := "UTC"
	if googleEvent.Start.TimeZone != "" {
		timezone = googleEvent.Start.TimeZone
	}

	event := &entities.Event{
		UserID:      userID,
		Title:       googleEvent.Summary,
		Description: googleEvent.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		Timezone:    timezone,
		Location:    googleEvent.Location,
		Attendees:   attendees,
		Status:      entities.EventStatusConfirmed,
		GoogleEventID: &googleEvent.ID,
	}

	return event, nil
}

// extractCalendarIDFromURI extracts calendar ID from Google Calendar resource URI
func (h *GoogleWebhookHandler) extractCalendarIDFromURI(resourceURI string) string {
	// Resource URI format: https://www.googleapis.com/calendar/v3/calendars/{calendarId}/events
	parts := strings.Split(resourceURI, "/")
	for i, part := range parts {
		if part == "calendars" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// SetupWebhookRequest represents the request to setup a webhook
type SetupWebhookRequest struct {
	CalendarID string `json:"calendar_id" binding:"required"`
}

// SetupWebhook sets up a webhook for a Google Calendar
func (h *GoogleWebhookHandler) SetupWebhook(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req SetupWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Google integration
	integration, err := h.googleIntegrationRepo.GetByUserID(c.Request.Context(), userID.(entities.UserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Google integration"})
		return
	}

	if integration == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Google integration not found"})
		return
	}

	// Get calendar sync configuration
	sync, err := h.googleCalendarSyncRepo.GetByCalendarID(c.Request.Context(), req.CalendarID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendar sync"})
		return
	}

	if sync == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar sync configuration not found"})
		return
	}

	// Setup webhook with Google Calendar API
	channelID := uuid.New().String()
	webhookURL := fmt.Sprintf("https://your-domain.com/api/v1/google/webhook") // TODO: Make this configurable
	
	accessToken := integration.AccessToken
	if integration.IsTokenExpiringSoon() {
		newTokens, err := h.oauth2Service.RefreshToken(c.Request.Context(), integration.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh Google token"})
			return
		}

		err = h.googleIntegrationRepo.UpdateTokens(c.Request.Context(), integration.ID, 
			newTokens.AccessToken, newTokens.RefreshToken, newTokens.Expiry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		accessToken = newTokens.AccessToken
	}

	err = h.calendarService.SetupWebhook(c.Request.Context(), accessToken, req.CalendarID, channelID, webhookURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup webhook with Google"})
		return
	}

	// Update sync configuration with webhook details
	sync.WebhookChannelID = channelID
	sync.WebhookURL = webhookURL
	sync.UpdatedAt = time.Now()

	err = h.googleCalendarSyncRepo.Update(c.Request.Context(), sync)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sync configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Webhook setup successfully",
		"channel_id": channelID,
		"webhook_url": webhookURL,
	})
}