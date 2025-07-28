package repositories

import (
	"context"
	"time"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type SyncConflictRepository interface {
	// Create a new conflict
	Create(ctx context.Context, conflict *entities.SyncConflict) error
	
	// Get conflict by ID
	GetByID(ctx context.Context, id string) (*entities.SyncConflict, error)
	
	// Get all pending conflicts for a user
	GetPendingByUserID(ctx context.Context, userID entities.UserID) ([]*entities.SyncConflict, error)
	
	// Get conflicts by calendar sync ID
	GetByCalendarSyncID(ctx context.Context, calendarSyncID string) ([]*entities.SyncConflict, error)
	
	// Update conflict status and resolution
	UpdateResolution(ctx context.Context, id string, resolution string, resolvedBy string, resolvedAt time.Time) error
	
	// Mark conflict as resolved
	MarkAsResolved(ctx context.Context, id string, resolvedBy string) error
	
	// Delete conflict
	Delete(ctx context.Context, id string) error
	
	// Get conflicts by status
	GetByStatus(ctx context.Context, userID entities.UserID, status string) ([]*entities.SyncConflict, error)
	
	// Check if similar conflict already exists
	FindSimilarConflict(ctx context.Context, userID entities.UserID, localEventID, googleEventID string) (*entities.SyncConflict, error)
	
	// Get conflict statistics
	GetConflictStats(ctx context.Context, userID entities.UserID, since time.Time) (map[string]int, error)
}