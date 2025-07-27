package repositories

import (
	"context"
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type EventRepository interface {
	// Create a new event
	Create(ctx context.Context, event *entities.Event) error
	
	// Get event by ID
	GetByID(ctx context.Context, id entities.EventID) (*entities.Event, error)
	
	// Get all events for a user
	GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Event, error)
	
	// Get events for a user within a time range
	GetByUserIDAndTimeRange(ctx context.Context, userID entities.UserID, start, end time.Time) ([]*entities.Event, error)
	
	// Alias for GetByUserIDAndTimeRange for consistency
	GetByTimeRange(ctx context.Context, userID entities.UserID, start, end time.Time) ([]*entities.Event, error)
	
	// Get events for a specific goal
	GetByGoalID(ctx context.Context, goalID entities.GoalID) ([]*entities.Event, error)
	
	// Get events by external ID (for sync purposes)
	GetByExternalID(ctx context.Context, userID entities.UserID, externalID string) (*entities.Event, error)
	
	// Get events by external source
	GetByExternalSource(ctx context.Context, userID entities.UserID, source string) ([]*entities.Event, error)
	
	// Get event by Google Event ID
	GetByGoogleEventID(ctx context.Context, googleEventID string) (*entities.Event, error)
	
	// Get upcoming events for a user
	GetUpcoming(ctx context.Context, userID entities.UserID, limit int) ([]*entities.Event, error)
	
	// Get events for today
	GetForToday(ctx context.Context, userID entities.UserID, timezone string) ([]*entities.Event, error)
	
	// Get recurring events
	GetRecurring(ctx context.Context, userID entities.UserID) ([]*entities.Event, error)
	
	// Get events by status
	GetByStatus(ctx context.Context, userID entities.UserID, status entities.EventStatus) ([]*entities.Event, error)
	
	// Update event
	Update(ctx context.Context, event *entities.Event) error
	
	// Delete event
	Delete(ctx context.Context, id entities.EventID) error
	
	// Bulk create events (for recurring events)
	BulkCreate(ctx context.Context, events []*entities.Event) error
	
	// Check for time conflicts
	HasConflict(ctx context.Context, userID entities.UserID, start, end time.Time, excludeEventID *entities.EventID) (bool, error)
	
	// Get events with pagination
	GetByUserIDPaginated(ctx context.Context, userID entities.UserID, offset, limit int) ([]*entities.Event, int64, error)
	
	// Search events by title or description
	Search(ctx context.Context, userID entities.UserID, query string, limit int) ([]*entities.Event, error)
}