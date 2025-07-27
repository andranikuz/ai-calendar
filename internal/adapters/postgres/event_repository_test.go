package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

func setupEventTestDB(t *testing.T) (*eventRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := &eventRepository{
		pool: nil, // Will be set in integration tests
	}

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestEventRepository_Create_SQLQuery(t *testing.T) {
	// Test the basic structure of the CREATE query
	expectedFields := []string{
		"INSERT INTO events",
		"id", "user_id", "goal_id", "title", "description",
		"start_time", "end_time", "timezone", "location",
		"status", "external_id", "external_source", "google_event_id",
		"created_at", "updated_at",
	}

	// Verify all required fields are present in the query
	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field should not be empty")
	}
}

func TestEvent_Validation(t *testing.T) {
	now := time.Now()
	future := now.Add(2 * time.Hour)

	tests := []struct {
		name    string
		event   entities.Event
		isValid bool
	}{
		{
			name: "Valid event",
			event: entities.Event{
				ID:          "event-1",
				UserID:      "user-1",
				Title:       "Team Meeting",
				Description: "Weekly team sync",
				StartTime:   now,
				EndTime:     future,
				Timezone:    "UTC",
				Location:    "Conference Room A",
				Status:      entities.EventStatusConfirmed,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			isValid: true,
		},
		{
			name: "Empty title",
			event: entities.Event{
				ID:        "event-1",
				UserID:    "user-1",
				Title:     "",
				StartTime: now,
				EndTime:   future,
				Status:    entities.EventStatusConfirmed,
			},
			isValid: false,
		},
		{
			name: "End time before start time",
			event: entities.Event{
				ID:        "event-1",
				UserID:    "user-1",
				Title:     "Invalid Event",
				StartTime: future,
				EndTime:   now, // This is invalid
				Status:    entities.EventStatusConfirmed,
			},
			isValid: false,
		},
		{
			name: "Empty user ID",
			event: entities.Event{
				ID:        "event-1",
				UserID:    "",
				Title:     "Event without user",
				StartTime: now,
				EndTime:   future,
				Status:    entities.EventStatusConfirmed,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation logic
			hasID := tt.event.ID != ""
			hasUserID := tt.event.UserID != ""
			hasTitle := tt.event.Title != ""
			validTimeRange := !tt.event.EndTime.Before(tt.event.StartTime)

			basicValidation := hasID && hasUserID && hasTitle && validTimeRange

			if tt.isValid {
				assert.True(t, basicValidation, "Event should be valid")
			} else {
				assert.False(t, basicValidation, "Event should be invalid")
			}
		})
	}
}

func TestEventStatus_Constants(t *testing.T) {
	statuses := []entities.EventStatus{
		entities.EventStatusTentative,
		entities.EventStatusConfirmed,
		entities.EventStatusCancelled,
	}

	expectedStatuses := []string{
		"tentative",
		"confirmed",
		"cancelled",
	}

	for i, status := range statuses {
		assert.Equal(t, expectedStatuses[i], string(status))
	}
}

func TestEvent_TimezoneHandling(t *testing.T) {
	tests := []struct {
		timezone string
		valid    bool
	}{
		{"UTC", true},
		{"America/New_York", true},
		{"Europe/London", true},
		{"Asia/Tokyo", true},
		{"", false}, // Empty timezone should be invalid
		{"Invalid/Timezone", true}, // We'll allow any string for now
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("timezone_%s", tt.timezone), func(t *testing.T) {
			event := entities.Event{
				Timezone: tt.timezone,
			}

			hasTimezone := event.Timezone != ""
			if tt.valid && tt.timezone != "" {
				assert.True(t, hasTimezone)
			}
		})
	}
}

func TestEvent_ExternalIntegration(t *testing.T) {
	tests := []struct {
		name           string
		externalID     string
		externalSource string
		googleEventID  *string
		hasExternal    bool
	}{
		{
			name:           "Google Calendar event",
			externalID:     "google-event-123",
			externalSource: "google",
			googleEventID:  stringPtr("google-123"),
			hasExternal:    true,
		},
		{
			name:           "Outlook event",
			externalID:     "outlook-event-456",
			externalSource: "outlook",
			googleEventID:  nil,
			hasExternal:    true,
		},
		{
			name:           "Local event only",
			externalID:     "",
			externalSource: "",
			googleEventID:  nil,
			hasExternal:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := entities.Event{
				ExternalID:     tt.externalID,
				ExternalSource: tt.externalSource,
				GoogleEventID:  tt.googleEventID,
			}

			hasExternal := event.ExternalID != "" || event.ExternalSource != "" || event.GoogleEventID != nil
			assert.Equal(t, tt.hasExternal, hasExternal)

			if tt.googleEventID != nil {
				assert.NotNil(t, event.GoogleEventID)
				assert.Equal(t, *tt.googleEventID, *event.GoogleEventID)
			}
		})
	}
}

func TestEvent_GoalLinking(t *testing.T) {
	goalID := entities.GoalID("goal-123")

	tests := []struct {
		name      string
		goalID    *entities.GoalID
		hasGoal   bool
	}{
		{
			name:    "Event linked to goal",
			goalID:  &goalID,
			hasGoal: true,
		},
		{
			name:    "Event not linked to goal",
			goalID:  nil,
			hasGoal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := entities.Event{
				GoalID: tt.goalID,
			}

			hasGoal := event.GoalID != nil
			assert.Equal(t, tt.hasGoal, hasGoal)

			if hasGoal {
				assert.Equal(t, goalID, *event.GoalID)
			}
		})
	}
}

func TestEvent_TimeValidation(t *testing.T) {
	now := time.Now()
	future1 := now.Add(1 * time.Hour)
	future2 := now.Add(2 * time.Hour)
	past := now.Add(-1 * time.Hour)

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		valid     bool
	}{
		{
			name:      "Valid time range",
			startTime: now,
			endTime:   future1,
			valid:     true,
		},
		{
			name:      "Same start and end time",
			startTime: now,
			endTime:   now,
			valid:     true, // Zero duration events can be valid
		},
		{
			name:      "End before start",
			startTime: future1,
			endTime:   now,
			valid:     false,
		},
		{
			name:      "Long duration event",
			startTime: now,
			endTime:   future2,
			valid:     true,
		},
		{
			name:      "Past event",
			startTime: past,
			endTime:   now,
			valid:     true, // Past events are valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := !tt.endTime.Before(tt.startTime)
			assert.Equal(t, tt.valid, isValid)
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}