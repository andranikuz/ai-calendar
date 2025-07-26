package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
)

type eventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) repositories.EventRepository {
	return &eventRepository{pool: pool}
}

func (r *eventRepository) Create(ctx context.Context, event *entities.Event) error {
	query := `
		INSERT INTO events (
			id, user_id, goal_id, title, description, start_time, end_time,
			timezone, recurrence, location, attendees, status, external_id,
			external_source, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err := r.pool.Exec(ctx, query,
		event.ID, event.UserID, event.GoalID, event.Title, event.Description,
		event.StartTime, event.EndTime, event.Timezone, event.Recurrence,
		event.Location, event.Attendees, event.Status, event.ExternalID,
		event.ExternalSource, event.CreatedAt, event.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	
	return nil
}

func (r *eventRepository) GetByID(ctx context.Context, id entities.EventID) (*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE id = $1`

	var event entities.Event
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
		&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
		&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
		&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}

	return &event, nil
}

func (r *eventRepository) GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1
		ORDER BY start_time ASC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by user ID: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *eventRepository) GetByUserIDAndTimeRange(ctx context.Context, userID entities.UserID, start, end time.Time) ([]*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
		ORDER BY start_time ASC`

	rows, err := r.pool.Query(ctx, query, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by time range: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

// GetByTimeRange is an alias for GetByUserIDAndTimeRange for consistency
func (r *eventRepository) GetByTimeRange(ctx context.Context, userID entities.UserID, start, end time.Time) ([]*entities.Event, error) {
	return r.GetByUserIDAndTimeRange(ctx, userID, start, end)
}

func (r *eventRepository) GetByGoalID(ctx context.Context, goalID entities.GoalID) ([]*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE goal_id = $1
		ORDER BY start_time ASC`

	rows, err := r.pool.Query(ctx, query, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by goal ID: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *eventRepository) GetByExternalID(ctx context.Context, userID entities.UserID, externalID string) (*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1 AND external_id = $2`

	var event entities.Event
	err := r.pool.QueryRow(ctx, query, userID, externalID).Scan(
		&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
		&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
		&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
		&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get event by external ID: %w", err)
	}

	return &event, nil
}

func (r *eventRepository) GetByExternalSource(ctx context.Context, userID entities.UserID, source string) ([]*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1 AND external_source = $2
		ORDER BY start_time ASC`

	rows, err := r.pool.Query(ctx, query, userID, source)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by external source: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *eventRepository) GetUpcoming(ctx context.Context, userID entities.UserID, limit int) ([]*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1 AND start_time > NOW() AND status != 'cancelled'
		ORDER BY start_time ASC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *eventRepository) GetForToday(ctx context.Context, userID entities.UserID, timezone string) ([]*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1 
		  AND DATE(start_time AT TIME ZONE $2) = CURRENT_DATE
		  AND status != 'cancelled'
		ORDER BY start_time ASC`

	rows, err := r.pool.Query(ctx, query, userID, timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's events: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *eventRepository) GetRecurring(ctx context.Context, userID entities.UserID) ([]*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1 AND recurrence IS NOT NULL AND recurrence != '{}'
		ORDER BY start_time ASC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recurring events: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *eventRepository) GetByStatus(ctx context.Context, userID entities.UserID, status entities.EventStatus) ([]*entities.Event, error) {
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1 AND status = $2
		ORDER BY start_time ASC`

	rows, err := r.pool.Query(ctx, query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by status: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *eventRepository) Update(ctx context.Context, event *entities.Event) error {
	query := `
		UPDATE events 
		SET goal_id = $2, title = $3, description = $4, start_time = $5, end_time = $6,
			timezone = $7, recurrence = $8, location = $9, attendees = $10, status = $11,
			external_id = $12, external_source = $13, updated_at = $14
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query,
		event.ID, event.GoalID, event.Title, event.Description, event.StartTime,
		event.EndTime, event.Timezone, event.Recurrence, event.Location,
		event.Attendees, event.Status, event.ExternalID, event.ExternalSource,
		event.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id entities.EventID) error {
	query := `DELETE FROM events WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func (r *eventRepository) BulkCreate(ctx context.Context, events []*entities.Event) error {
	if len(events) == 0 {
		return nil
	}

	query := `
		INSERT INTO events (
			id, user_id, goal_id, title, description, start_time, end_time,
			timezone, recurrence, location, attendees, status, external_id,
			external_source, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	batch := &pgx.Batch{}
	for _, event := range events {
		batch.Queue(query,
			event.ID, event.UserID, event.GoalID, event.Title, event.Description,
			event.StartTime, event.EndTime, event.Timezone, event.Recurrence,
			event.Location, event.Attendees, event.Status, event.ExternalID,
			event.ExternalSource, event.CreatedAt, event.UpdatedAt,
		)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(events); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to bulk create events: %w", err)
		}
	}

	return nil
}

func (r *eventRepository) HasConflict(ctx context.Context, userID entities.UserID, start, end time.Time, excludeEventID *entities.EventID) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM events 
		WHERE user_id = $1 
		  AND status != 'cancelled'
		  AND (start_time < $3 AND end_time > $2)`

	args := []interface{}{userID, start, end}

	if excludeEventID != nil {
		query += " AND id != $4"
		args = append(args, *excludeEventID)
	}

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for conflicts: %w", err)
	}

	return count > 0, nil
}

func (r *eventRepository) GetByUserIDPaginated(ctx context.Context, userID entities.UserID, offset, limit int) ([]*entities.Event, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM events WHERE user_id = $1`
	var totalCount int64
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get events count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1
		ORDER BY start_time ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated events: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, totalCount, nil
}

func (r *eventRepository) Search(ctx context.Context, userID entities.UserID, query string, limit int) ([]*entities.Event, error) {
	searchQuery := `
		SELECT id, user_id, goal_id, title, description, start_time, end_time,
			   timezone, recurrence, location, attendees, status, external_id,
			   external_source, created_at, updated_at
		FROM events 
		WHERE user_id = $1 
		  AND (title ILIKE $2 OR description ILIKE $2 OR location ILIKE $2)
		ORDER BY start_time ASC
		LIMIT $3`

	searchPattern := "%" + query + "%"
	rows, err := r.pool.Query(ctx, searchQuery, userID, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search events: %w", err)
	}
	defer rows.Close()

	var events []*entities.Event
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.GoalID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.Timezone, &event.Recurrence,
			&event.Location, &event.Attendees, &event.Status, &event.ExternalID,
			&event.ExternalSource, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}