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

type goalRepository struct {
	pool *pgxpool.Pool
}

func NewGoalRepository(pool *pgxpool.Pool) repositories.GoalRepository {
	return &goalRepository{pool: pool}
}

func (r *goalRepository) Create(ctx context.Context, goal *entities.Goal) error {
	query := `
		INSERT INTO goals (
			id, user_id, title, description, category, priority, status, 
			progress, deadline, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.pool.Exec(ctx, query,
		goal.ID, goal.UserID, goal.Title, goal.Description, 
		goal.Category, goal.Priority, goal.Status, goal.Progress,
		goal.Deadline, goal.CreatedAt, goal.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create goal: %w", err)
	}
	
	return nil
}

func (r *goalRepository) GetByID(ctx context.Context, id entities.GoalID) (*entities.Goal, error) {
	query := `
		SELECT id, user_id, title, description, category, priority, status, 
			   progress, deadline, created_at, updated_at
		FROM goals 
		WHERE id = $1`

	var goal entities.Goal
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&goal.ID, &goal.UserID, &goal.Title, &goal.Description,
		&goal.Category, &goal.Priority, &goal.Status, &goal.Progress,
		&goal.Deadline, &goal.CreatedAt, &goal.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get goal by ID: %w", err)
	}

	return &goal, nil
}

func (r *goalRepository) GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
	query := `
		SELECT id, user_id, title, description, category, priority, status, 
			   progress, deadline, created_at, updated_at
		FROM goals 
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals by user ID: %w", err)
	}
	defer rows.Close()

	var goals []*entities.Goal
	for rows.Next() {
		var goal entities.Goal
		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Title, &goal.Description,
			&goal.Category, &goal.Priority, &goal.Status, &goal.Progress,
			&goal.Deadline, &goal.CreatedAt, &goal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, &goal)
	}

	return goals, nil
}

func (r *goalRepository) GetByUserIDAndStatus(ctx context.Context, userID entities.UserID, status entities.GoalStatus) ([]*entities.Goal, error) {
	query := `
		SELECT id, user_id, title, description, category, priority, status, 
			   progress, deadline, created_at, updated_at
		FROM goals 
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals by status: %w", err)
	}
	defer rows.Close()

	var goals []*entities.Goal
	for rows.Next() {
		var goal entities.Goal
		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Title, &goal.Description,
			&goal.Category, &goal.Priority, &goal.Status, &goal.Progress,
			&goal.Deadline, &goal.CreatedAt, &goal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, &goal)
	}

	return goals, nil
}

func (r *goalRepository) GetByUserIDAndCategory(ctx context.Context, userID entities.UserID, category entities.GoalCategory) ([]*entities.Goal, error) {
	query := `
		SELECT id, user_id, title, description, category, priority, status, 
			   progress, deadline, created_at, updated_at
		FROM goals 
		WHERE user_id = $1 AND category = $2
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, userID, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals by category: %w", err)
	}
	defer rows.Close()

	var goals []*entities.Goal
	for rows.Next() {
		var goal entities.Goal
		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Title, &goal.Description,
			&goal.Category, &goal.Priority, &goal.Status, &goal.Progress,
			&goal.Deadline, &goal.CreatedAt, &goal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, &goal)
	}

	return goals, nil
}

func (r *goalRepository) GetByDeadlineBefore(ctx context.Context, userID entities.UserID, deadline time.Time) ([]*entities.Goal, error) {
	query := `
		SELECT id, user_id, title, description, category, priority, status, 
			   progress, deadline, created_at, updated_at
		FROM goals 
		WHERE user_id = $1 AND deadline <= $2 AND status != 'completed'
		ORDER BY deadline ASC`

	rows, err := r.pool.Query(ctx, query, userID, deadline)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals by deadline: %w", err)
	}
	defer rows.Close()

	var goals []*entities.Goal
	for rows.Next() {
		var goal entities.Goal
		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Title, &goal.Description,
			&goal.Category, &goal.Priority, &goal.Status, &goal.Progress,
			&goal.Deadline, &goal.CreatedAt, &goal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, &goal)
	}

	return goals, nil
}

func (r *goalRepository) Update(ctx context.Context, goal *entities.Goal) error {
	query := `
		UPDATE goals 
		SET title = $2, description = $3, category = $4, priority = $5, 
			status = $6, progress = $7, deadline = $8, updated_at = $9
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query,
		goal.ID, goal.Title, goal.Description, goal.Category,
		goal.Priority, goal.Status, goal.Progress, goal.Deadline,
		goal.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}

	return nil
}

func (r *goalRepository) Delete(ctx context.Context, id entities.GoalID) error {
	query := `DELETE FROM goals WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}

	return nil
}

func (r *goalRepository) UpdateProgress(ctx context.Context, goalID entities.GoalID, progress int) error {
	query := `
		UPDATE goals 
		SET progress = $2, updated_at = $3
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, goalID, progress, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update goal progress: %w", err)
	}

	return nil
}

func (r *goalRepository) GetByUserIDPaginated(ctx context.Context, userID entities.UserID, offset, limit int) ([]*entities.Goal, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM goals WHERE user_id = $1`
	var totalCount int64
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get goals count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, user_id, title, description, category, priority, status, 
			   progress, deadline, created_at, updated_at
		FROM goals 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated goals: %w", err)
	}
	defer rows.Close()

	var goals []*entities.Goal
	for rows.Next() {
		var goal entities.Goal
		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Title, &goal.Description,
			&goal.Category, &goal.Priority, &goal.Status, &goal.Progress,
			&goal.Deadline, &goal.CreatedAt, &goal.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, &goal)
	}

	return goals, totalCount, nil
}