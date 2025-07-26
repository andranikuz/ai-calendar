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
	LastSyncAt          *time.Time              `json:"last_sync_at"`
	LastSyncError       string                  `json:"last_sync_error"`
	SyncToken           string                  `json:"sync_token"`
	Settings            CalendarSyncSettings    `json:"settings"`
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

type CalendarSyncSettings struct {
	SyncInterval      time.Duration `json:"sync_interval"`       // How often to sync
	AutoSync          bool          `json:"auto_sync"`           // Enable automatic sync
	SyncPastEvents    bool          `json:"sync_past_events"`    // Sync events from the past
	SyncFutureEvents  bool          `json:"sync_future_events"`  // Sync future events
	ConflictResolution string       `json:"conflict_resolution"` // "google_wins", "local_wins", "manual"
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

	if gcs.LastSyncAt == nil {
		return true
	}

	return time.Now().Sub(*gcs.LastSyncAt) >= gcs.Settings.SyncInterval
}

// Default settings
func DefaultCalendarSyncSettings() CalendarSyncSettings {
	return CalendarSyncSettings{
		SyncInterval:       15 * time.Minute,
		AutoSync:           true,
		SyncPastEvents:     false,
		SyncFutureEvents:   true,
		ConflictResolution: "google_wins",
	}
}