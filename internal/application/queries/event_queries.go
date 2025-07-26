package queries

import (
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

// Event queries

type GetEventByIDQuery struct {
	EventID entities.EventID `json:"event_id"`
	UserID  entities.UserID  `json:"user_id"`
}

type GetEventResult struct {
	Event *entities.Event `json:"event"`
}

type GetEventsByUserIDQuery struct {
	UserID entities.UserID `json:"user_id"`
	Offset int             `json:"offset"`
	Limit  int             `json:"limit"`
}

type GetEventsResult struct {
	Events     []*entities.Event `json:"events"`
	TotalCount int64             `json:"total_count"`
	Offset     int               `json:"offset"`
	Limit      int               `json:"limit"`
}

type GetEventsByTimeRangeQuery struct {
	UserID    entities.UserID `json:"user_id"`
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time"`
}

type GetEventsByTimeRangeResult struct {
	Events []*entities.Event `json:"events"`
}

type GetEventsByGoalIDQuery struct {
	GoalID entities.GoalID `json:"goal_id"`
	UserID entities.UserID `json:"user_id"`
}

type GetEventsByGoalIDResult struct {
	Events []*entities.Event `json:"events"`
}

type GetUpcomingEventsQuery struct {
	UserID entities.UserID `json:"user_id"`
	Limit  int             `json:"limit"`
}

type GetUpcomingEventsResult struct {
	Events []*entities.Event `json:"events"`
}

type GetTodayEventsQuery struct {
	UserID   entities.UserID `json:"user_id"`
	Timezone string          `json:"timezone"`
}

type GetTodayEventsResult struct {
	Events []*entities.Event `json:"events"`
}

type GetRecurringEventsQuery struct {
	UserID entities.UserID `json:"user_id"`
}

type GetRecurringEventsResult struct {
	Events []*entities.Event `json:"events"`
}

type GetEventsByStatusQuery struct {
	UserID entities.UserID     `json:"user_id"`
	Status entities.EventStatus `json:"status"`
	Offset int                 `json:"offset"`
	Limit  int                 `json:"limit"`
}

type GetEventsByStatusResult struct {
	Events     []*entities.Event `json:"events"`
	TotalCount int64             `json:"total_count"`
}

type SearchEventsQuery struct {
	UserID entities.UserID `json:"user_id"`
	Query  string          `json:"query"`
	Limit  int             `json:"limit"`
}

type SearchEventsResult struct {
	Events []*entities.Event `json:"events"`
}

type CheckEventConflictQuery struct {
	UserID          entities.UserID   `json:"user_id"`
	StartTime       time.Time         `json:"start_time"`
	EndTime         time.Time         `json:"end_time"`
	ExcludeEventID  *entities.EventID `json:"exclude_event_id,omitempty"`
}

type CheckEventConflictResult struct {
	HasConflict bool `json:"has_conflict"`
}

type GetEventsByExternalSourceQuery struct {
	UserID         entities.UserID `json:"user_id"`
	ExternalSource string          `json:"external_source"`
}

type GetEventsByExternalSourceResult struct {
	Events []*entities.Event `json:"events"`
}

type GetEventStatsQuery struct {
	UserID    entities.UserID `json:"user_id"`
	StartDate time.Time       `json:"start_date"`
	EndDate   time.Time       `json:"end_date"`
}

type GetEventStatsResult struct {
	TotalEvents      int                              `json:"total_events"`
	EventsByStatus   map[entities.EventStatus]int     `json:"events_by_status"`
	EventsByCategory map[entities.GoalCategory]int    `json:"events_by_category"`
	TotalHours       float64                          `json:"total_hours"`
	AveragePerDay    float64                          `json:"average_per_day"`
}