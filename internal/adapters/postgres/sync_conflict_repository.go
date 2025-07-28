package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
)

type syncConflictRepository struct {
	db *sql.DB
}

func NewSyncConflictRepository(db *sql.DB) repositories.SyncConflictRepository {
	return &syncConflictRepository{db: db}
}

func (r *syncConflictRepository) Create(ctx context.Context, conflict *entities.SyncConflict) error {
	if conflict.ID == "" {
		conflict.ID = uuid.New().String()
	}
	
	now := time.Now()
	conflict.CreatedAt = now
	conflict.UpdatedAt = now
	
	if conflict.Status == "" {
		conflict.Status = "pending"
	}

	localEventJSON, err := json.Marshal(conflict.LocalEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal local event: %w", err)
	}

	googleEventJSON, err := json.Marshal(conflict.GoogleEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal google event: %w", err)
	}

	query := `
		INSERT INTO sync_conflicts (
			id, user_id, calendar_sync_id, conflict_type, local_event, google_event,
			description, resolution, resolved_at, resolved_by, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = r.db.ExecContext(ctx, query,
		conflict.ID, conflict.UserID, conflict.CalendarSyncID, conflict.ConflictType,
		localEventJSON, googleEventJSON, conflict.Description, conflict.Resolution,
		conflict.ResolvedAt, conflict.ResolvedBy, conflict.Status,
		conflict.CreatedAt, conflict.UpdatedAt,
	)

	return err
}

func (r *syncConflictRepository) GetByID(ctx context.Context, id string) (*entities.SyncConflict, error) {
	query := `
		SELECT id, user_id, calendar_sync_id, conflict_type, local_event, google_event,
			   description, resolution, resolved_at, resolved_by, status, created_at, updated_at
		FROM sync_conflicts WHERE id = $1`

	var conflict entities.SyncConflict
	var localEventJSON, googleEventJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conflict.ID, &conflict.UserID, &conflict.CalendarSyncID, &conflict.ConflictType,
		&localEventJSON, &googleEventJSON, &conflict.Description, &conflict.Resolution,
		&conflict.ResolvedAt, &conflict.ResolvedBy, &conflict.Status,
		&conflict.CreatedAt, &conflict.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sync conflict not found")
		}
		return nil, err
	}

	if len(localEventJSON) > 0 {
		if err := json.Unmarshal(localEventJSON, &conflict.LocalEvent); err != nil {
			return nil, fmt.Errorf("failed to unmarshal local event: %w", err)
		}
	}

	if len(googleEventJSON) > 0 {
		if err := json.Unmarshal(googleEventJSON, &conflict.GoogleEvent); err != nil {
			return nil, fmt.Errorf("failed to unmarshal google event: %w", err)
		}
	}

	return &conflict, nil
}

func (r *syncConflictRepository) GetPendingByUserID(ctx context.Context, userID entities.UserID) ([]*entities.SyncConflict, error) {
	return r.GetByStatus(ctx, userID, "pending")
}

func (r *syncConflictRepository) GetByCalendarSyncID(ctx context.Context, calendarSyncID string) ([]*entities.SyncConflict, error) {
	query := `
		SELECT id, user_id, calendar_sync_id, conflict_type, local_event, google_event,
			   description, resolution, resolved_at, resolved_by, status, created_at, updated_at
		FROM sync_conflicts WHERE calendar_sync_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, calendarSyncID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanConflicts(rows)
}

func (r *syncConflictRepository) UpdateResolution(ctx context.Context, id string, resolution string, resolvedBy string, resolvedAt time.Time) error {
	query := `
		UPDATE sync_conflicts 
		SET resolution = $2, resolved_by = $3, resolved_at = $4, status = 'resolved', updated_at = $5
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, resolution, resolvedBy, resolvedAt, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sync conflict not found")
	}

	return nil
}

func (r *syncConflictRepository) MarkAsResolved(ctx context.Context, id string, resolvedBy string) error {
	return r.UpdateResolution(ctx, id, "resolved", resolvedBy, time.Now())
}

func (r *syncConflictRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sync_conflicts WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sync conflict not found")
	}

	return nil
}

func (r *syncConflictRepository) GetByStatus(ctx context.Context, userID entities.UserID, status string) ([]*entities.SyncConflict, error) {
	query := `
		SELECT id, user_id, calendar_sync_id, conflict_type, local_event, google_event,
			   description, resolution, resolved_at, resolved_by, status, created_at, updated_at
		FROM sync_conflicts WHERE user_id = $1 AND status = $2 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanConflicts(rows)
}

func (r *syncConflictRepository) FindSimilarConflict(ctx context.Context, userID entities.UserID, localEventID, googleEventID string) (*entities.SyncConflict, error) {
	query := `
		SELECT id, user_id, calendar_sync_id, conflict_type, local_event, google_event,
			   description, resolution, resolved_at, resolved_by, status, created_at, updated_at
		FROM sync_conflicts 
		WHERE user_id = $1 AND status = 'pending'
		  AND (local_event->>'id' = $2 OR google_event->>'id' = $3)
		ORDER BY created_at DESC
		LIMIT 1`

	var conflict entities.SyncConflict
	var localEventJSON, googleEventJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID, localEventID, googleEventID).Scan(
		&conflict.ID, &conflict.UserID, &conflict.CalendarSyncID, &conflict.ConflictType,
		&localEventJSON, &googleEventJSON, &conflict.Description, &conflict.Resolution,
		&conflict.ResolvedAt, &conflict.ResolvedBy, &conflict.Status,
		&conflict.CreatedAt, &conflict.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No similar conflict found
		}
		return nil, err
	}

	if len(localEventJSON) > 0 {
		if err := json.Unmarshal(localEventJSON, &conflict.LocalEvent); err != nil {
			return nil, fmt.Errorf("failed to unmarshal local event: %w", err)
		}
	}

	if len(googleEventJSON) > 0 {
		if err := json.Unmarshal(googleEventJSON, &conflict.GoogleEvent); err != nil {
			return nil, fmt.Errorf("failed to unmarshal google event: %w", err)
		}
	}

	return &conflict, nil
}

func (r *syncConflictRepository) GetConflictStats(ctx context.Context, userID entities.UserID, since time.Time) (map[string]int, error) {
	query := `
		SELECT conflict_type, COUNT(*) as count
		FROM sync_conflicts 
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY conflict_type`

	rows, err := r.db.QueryContext(ctx, query, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var conflictType string
		var count int
		if err := rows.Scan(&conflictType, &count); err != nil {
			return nil, err
		}
		stats[conflictType] = count
	}

	return stats, rows.Err()
}

func (r *syncConflictRepository) scanConflicts(rows *sql.Rows) ([]*entities.SyncConflict, error) {
	var conflicts []*entities.SyncConflict

	for rows.Next() {
		var conflict entities.SyncConflict
		var localEventJSON, googleEventJSON []byte

		err := rows.Scan(
			&conflict.ID, &conflict.UserID, &conflict.CalendarSyncID, &conflict.ConflictType,
			&localEventJSON, &googleEventJSON, &conflict.Description, &conflict.Resolution,
			&conflict.ResolvedAt, &conflict.ResolvedBy, &conflict.Status,
			&conflict.CreatedAt, &conflict.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if len(localEventJSON) > 0 {
			if err := json.Unmarshal(localEventJSON, &conflict.LocalEvent); err != nil {
				return nil, fmt.Errorf("failed to unmarshal local event: %w", err)
			}
		}

		if len(googleEventJSON) > 0 {
			if err := json.Unmarshal(googleEventJSON, &conflict.GoogleEvent); err != nil {
				return nil, fmt.Errorf("failed to unmarshal google event: %w", err)
			}
		}

		conflicts = append(conflicts, &conflict)
	}

	return conflicts, rows.Err()
}