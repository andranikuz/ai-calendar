package commands

import (
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

// Event commands

type CreateEventCommand struct {
	UserID         entities.UserID         `json:"user_id"`
	GoalID         *entities.GoalID        `json:"goal_id,omitempty"`
	Title          string                  `json:"title"`
	Description    string                  `json:"description"`
	StartTime      time.Time               `json:"start_time"`
	EndTime        time.Time               `json:"end_time"`
	Timezone       string                  `json:"timezone"`
	Recurrence     *entities.RecurrenceRule `json:"recurrence,omitempty"`
	Location       string                  `json:"location"`
	Attendees      []entities.Attendee     `json:"attendees"`
	Status         entities.EventStatus    `json:"status"`
	ExternalID     string                  `json:"external_id"`
	ExternalSource string                  `json:"external_source"`
}

type CreateEventResult struct {
	EventID   entities.EventID `json:"event_id"`
	CreatedAt time.Time        `json:"created_at"`
}

type UpdateEventCommand struct {
	EventID        entities.EventID         `json:"event_id"`
	UserID         entities.UserID          `json:"user_id"`
	GoalID         *entities.GoalID         `json:"goal_id,omitempty"`
	Title          *string                  `json:"title,omitempty"`
	Description    *string                  `json:"description,omitempty"`
	StartTime      *time.Time               `json:"start_time,omitempty"`
	EndTime        *time.Time               `json:"end_time,omitempty"`
	Timezone       *string                  `json:"timezone,omitempty"`
	Recurrence     *entities.RecurrenceRule `json:"recurrence,omitempty"`
	Location       *string                  `json:"location,omitempty"`
	Attendees      *[]entities.Attendee     `json:"attendees,omitempty"`
	Status         *entities.EventStatus    `json:"status,omitempty"`
	ExternalID     *string                  `json:"external_id,omitempty"`
	ExternalSource *string                  `json:"external_source,omitempty"`
}

type UpdateEventResult struct {
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteEventCommand struct {
	EventID entities.EventID `json:"event_id"`
	UserID  entities.UserID  `json:"user_id"`
}

type MoveEventCommand struct {
	EventID   entities.EventID `json:"event_id"`
	UserID    entities.UserID  `json:"user_id"`
	StartTime time.Time        `json:"start_time"`
	EndTime   time.Time        `json:"end_time"`
}

type MoveEventResult struct {
	UpdatedAt time.Time `json:"updated_at"`
}

type DuplicateEventCommand struct {
	EventID   entities.EventID `json:"event_id"`
	UserID    entities.UserID  `json:"user_id"`
	StartTime time.Time        `json:"start_time"`
	EndTime   time.Time        `json:"end_time"`
}

type DuplicateEventResult struct {
	EventID   entities.EventID `json:"event_id"`
	CreatedAt time.Time        `json:"created_at"`
}

type ChangeEventStatusCommand struct {
	EventID entities.EventID    `json:"event_id"`
	UserID  entities.UserID     `json:"user_id"`
	Status  entities.EventStatus `json:"status"`
}

type ChangeEventStatusResult struct {
	UpdatedAt time.Time `json:"updated_at"`
}

type LinkEventToGoalCommand struct {
	EventID entities.EventID `json:"event_id"`
	UserID  entities.UserID  `json:"user_id"`
	GoalID  entities.GoalID  `json:"goal_id"`
}

type LinkEventToGoalResult struct {
	UpdatedAt time.Time `json:"updated_at"`
}

type UnlinkEventFromGoalCommand struct {
	EventID entities.EventID `json:"event_id"`
	UserID  entities.UserID  `json:"user_id"`
}

type UnlinkEventFromGoalResult struct {
	UpdatedAt time.Time `json:"updated_at"`
}

type SyncEventCommand struct {
	UserID         entities.UserID  `json:"user_id"`
	ExternalSource string           `json:"external_source"`
	Events         []*entities.Event `json:"events"`
}

type SyncEventResult struct {
	CreatedCount int       `json:"created_count"`
	UpdatedCount int       `json:"updated_count"`
	DeletedCount int       `json:"deleted_count"`
	SyncedAt     time.Time `json:"synced_at"`
}