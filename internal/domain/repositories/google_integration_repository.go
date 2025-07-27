package repositories

import (
	"context"
	"time"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type GoogleIntegrationRepository interface {
	// Create a new Google integration
	Create(ctx context.Context, integration *entities.GoogleIntegration) error
	
	// Get integration by ID
	GetByID(ctx context.Context, id entities.GoogleIntegrationID) (*entities.GoogleIntegration, error)
	
	// Get integration by user ID
	GetByUserID(ctx context.Context, userID entities.UserID) (*entities.GoogleIntegration, error)
	
	// Get integration by Google user ID
	GetByGoogleUserID(ctx context.Context, googleUserID string) (*entities.GoogleIntegration, error)
	
	// Update integration
	Update(ctx context.Context, integration *entities.GoogleIntegration) error
	
	// Update tokens
	UpdateTokens(ctx context.Context, id entities.GoogleIntegrationID, accessToken, refreshToken string, expiresAt time.Time) error
	
	// Delete integration
	Delete(ctx context.Context, id entities.GoogleIntegrationID) error
	
	// Get integrations that need token refresh
	GetExpiringSoon(ctx context.Context, beforeTime time.Time) ([]*entities.GoogleIntegration, error)
	
	// Get all active integrations
	GetActive(ctx context.Context) ([]*entities.GoogleIntegration, error)
}

type GoogleCalendarSyncRepository interface {
	// Create a new calendar sync configuration
	Create(ctx context.Context, sync *entities.GoogleCalendarSync) error
	
	// Get sync configuration by ID
	GetByID(ctx context.Context, id string) (*entities.GoogleCalendarSync, error)
	
	// Get sync configurations by user ID
	GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.GoogleCalendarSync, error)
	
	// Get sync configurations by integration ID
	GetByIntegrationID(ctx context.Context, integrationID entities.GoogleIntegrationID) ([]*entities.GoogleCalendarSync, error)
	
	// Get sync configuration by calendar ID
	GetByCalendarID(ctx context.Context, calendarID string) (*entities.GoogleCalendarSync, error)
	
	// Get sync configuration by webhook channel ID
	GetByChannelID(ctx context.Context, channelID string) (*entities.GoogleCalendarSync, error)
	
	// Update sync configuration
	Update(ctx context.Context, sync *entities.GoogleCalendarSync) error
	
	// Update sync status
	UpdateSyncStatus(ctx context.Context, id string, status entities.CalendarSyncStatus, lastSyncAt *time.Time, syncError string) error
	
	// Update sync token
	UpdateSyncToken(ctx context.Context, id string, syncToken string) error
	
	// Delete sync configuration
	Delete(ctx context.Context, id string) error
	
	// Get configurations that need sync
	GetNeedingSync(ctx context.Context) ([]*entities.GoogleCalendarSync, error)
	
	// Get active sync configurations
	GetActive(ctx context.Context) ([]*entities.GoogleCalendarSync, error)
}