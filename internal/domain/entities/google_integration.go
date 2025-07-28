package entities

import (
	"time"
)

type GoogleIntegrationID string

type GoogleIntegration struct {
	ID           GoogleIntegrationID `json:"id"`
	UserID       UserID              `json:"user_id"`
	GoogleUserID string              `json:"google_user_id"`
	Email        string              `json:"email"`
	Name         string              `json:"name"`
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	TokenType    string              `json:"token_type"`
	ExpiresAt    time.Time           `json:"expires_at"`
	Scopes       []string            `json:"scopes"`
	CalendarID   string              `json:"calendar_id"` // Primary calendar ID
	Enabled      bool                `json:"enabled"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type GoogleCalendarSync struct {
	ID                  string                  `json:"id"`
	UserID              UserID                  `json:"user_id"`
	GoogleIntegrationID GoogleIntegrationID     `json:"google_integration_id"`
	CalendarID          string                  `json:"calendar_id"`
	CalendarName        string                  `json:"calendar_name"`
	SyncDirection       CalendarSyncDirection   `json:"sync_direction"`
	SyncStatus          CalendarSyncStatus      `json:"sync_status"`
	LastSyncAt          time.Time               `json:"last_sync_at"`
	LastSyncError       string                  `json:"last_sync_error"`
	SyncToken           string                  `json:"sync_token"`
	Settings            CalendarSyncSettings    `json:"settings"`
	// Webhook fields
	WebhookChannelID    string                  `json:"webhook_channel_id"`
	WebhookURL          string                  `json:"webhook_url"`
	WebhookResourceID   string                  `json:"webhook_resource_id"`
	WebhookExpiresAt    *time.Time              `json:"webhook_expires_at"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
}

type CalendarSyncDirection string

const (
	SyncDirectionBidirectional CalendarSyncDirection = "bidirectional"
	SyncDirectionFromGoogle    CalendarSyncDirection = "from_google"
	SyncDirectionToGoogle      CalendarSyncDirection = "to_google"
)

type CalendarSyncStatus string

const (
	SyncStatusActive    CalendarSyncStatus = "active"
	SyncStatusPaused    CalendarSyncStatus = "paused"
	SyncStatusError     CalendarSyncStatus = "error"
	SyncStatusDisabled  CalendarSyncStatus = "disabled"
)

type ConflictResolutionStrategy string

const (
	ConflictResolutionGoogleWins ConflictResolutionStrategy = "google_wins"
	ConflictResolutionLocalWins  ConflictResolutionStrategy = "local_wins"
	ConflictResolutionManual     ConflictResolutionStrategy = "manual"
)

type CalendarSyncSettings struct {
	SyncInterval       time.Duration              `json:"sync_interval"`        // How often to sync
	AutoSync           bool                       `json:"auto_sync"`            // Enable automatic sync
	SyncPastEvents     bool                       `json:"sync_past_events"`     // Sync events from the past
	SyncFutureEvents   bool                       `json:"sync_future_events"`   // Sync future events
	ConflictResolution ConflictResolutionStrategy `json:"conflict_resolution"`  // Strategy for handling conflicts
}

// ConflictType represents the type of synchronization conflict
type ConflictType string

const (
	ConflictTypeTimeOverlap    ConflictType = "time_overlap"     // Events overlap in time
	ConflictTypeContentDiff    ConflictType = "content_diff"     // Same event, different content
	ConflictTypeDuplicateEvent ConflictType = "duplicate_event"  // Duplicate events found
	ConflictTypeDeletedEvent   ConflictType = "deleted_event"    // Event deleted in one source
)

// SyncConflict represents a conflict that occurred during synchronization
type SyncConflict struct {
	ID              string        `json:"id"`
	UserID          UserID        `json:"user_id"`
	CalendarSyncID  string        `json:"calendar_sync_id"`
	ConflictType    ConflictType  `json:"conflict_type"`
	LocalEvent      *Event        `json:"local_event"`
	GoogleEvent     *Event        `json:"google_event"`
	Description     string        `json:"description"`
	Resolution      string        `json:"resolution"`      // How it was resolved
	ResolvedAt      *time.Time    `json:"resolved_at"`
	ResolvedBy      string        `json:"resolved_by"`     // "auto" or user ID
	Status          string        `json:"status"`          // "pending", "resolved", "ignored"
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// ConflictResolutionAction represents an action to resolve a conflict
type ConflictResolutionAction struct {
	Action      string            `json:"action"`       // "use_local", "use_google", "merge", "ignore"
	EventData   map[string]interface{} `json:"event_data"`   // Event data to use for resolution
	Resolution  string            `json:"resolution"`   // Human-readable description
}

// Validation methods
func (gi *GoogleIntegration) IsValid() bool {
	return gi.UserID != "" && 
		   gi.GoogleUserID != "" &&
		   gi.Email != "" &&
		   gi.AccessToken != ""
}

func (gi *GoogleIntegration) IsTokenExpired() bool {
	return time.Now().After(gi.ExpiresAt)
}

func (gi *GoogleIntegration) IsTokenExpiringSoon() bool {
	return time.Now().Add(5 * time.Minute).After(gi.ExpiresAt)
}

func (gcs *GoogleCalendarSync) IsValid() bool {
	return gcs.UserID != "" &&
		   gcs.GoogleIntegrationID != "" &&
		   gcs.CalendarID != ""
}

func (gcs *GoogleCalendarSync) IsActive() bool {
	return gcs.SyncStatus == SyncStatusActive
}

func (gcs *GoogleCalendarSync) NeedsSync() bool {
	if !gcs.IsActive() {
		return false
	}

	if gcs.LastSyncAt.IsZero() {
		return true
	}

	return time.Now().Sub(gcs.LastSyncAt) >= gcs.Settings.SyncInterval
}

// Default settings
func DefaultCalendarSyncSettings() CalendarSyncSettings {
	return CalendarSyncSettings{
		SyncInterval:       15 * time.Minute,
		AutoSync:           true,
		SyncPastEvents:     false,
		SyncFutureEvents:   true,
		ConflictResolution: ConflictResolutionGoogleWins,
	}
}