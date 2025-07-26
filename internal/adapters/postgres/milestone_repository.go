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

type milestoneRepository struct {
	pool *pgxpool.Pool
}

func NewMilestoneRepository(pool *pgxpool.Pool) repositories.MilestoneRepository {
	return &milestoneRepository{pool: pool}
}

func (r *milestoneRepository) Create(ctx context.Context, milestone *entities.Milestone) error {
	query := `
		INSERT INTO milestones (
			id, goal_id, title, description, target_date,
			completed, completed_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.pool.Exec(ctx, query,
		milestone.ID, milestone.GoalID, milestone.Title, milestone.Description,
		milestone.TargetDate, milestone.Completed, milestone.CompletedAt,
		milestone.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}
	
	return nil
}

func (r *milestoneRepository) GetByID(ctx context.Context, id entities.MilestoneID) (*entities.Milestone, error) {
	query := `
		SELECT id, goal_id, title, description, target_date,
			   completed, completed_at, created_at
		FROM milestones 
		WHERE id = $1`

	var milestone entities.Milestone
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&milestone.ID, &milestone.GoalID, &milestone.Title, &milestone.Description,
		&milestone.TargetDate, &milestone.Completed, &milestone.CompletedAt,
		&milestone.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get milestone by ID: %w", err)
	}

	return &milestone, nil
}

func (r *milestoneRepository) GetByGoalID(ctx context.Context, goalID entities.GoalID) ([]*entities.Milestone, error) {
	query := `
		SELECT id, goal_id, title, description, target_date,
			   completed, completed_at, created_at
		FROM milestones 
		WHERE goal_id = $1
		ORDER BY target_date ASC`

	rows, err := r.pool.Query(ctx, query, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones by goal ID: %w", err)
	}
	defer rows.Close()

	var milestones []*entities.Milestone
	for rows.Next() {
		var milestone entities.Milestone
		err := rows.Scan(
			&milestone.ID, &milestone.GoalID, &milestone.Title, &milestone.Description,
			&milestone.TargetDate, &milestone.Completed, &milestone.CompletedAt,
			&milestone.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}
		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

func (r *milestoneRepository) GetCompletedByGoalID(ctx context.Context, goalID entities.GoalID) ([]*entities.Milestone, error) {
	query := `
		SELECT id, goal_id, title, description, target_date,
			   completed, completed_at, created_at
		FROM milestones 
		WHERE goal_id = $1 AND completed = true
		ORDER BY completed_at DESC`

	rows, err := r.pool.Query(ctx, query, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed milestones: %w", err)
	}
	defer rows.Close()

	var milestones []*entities.Milestone
	for rows.Next() {
		var milestone entities.Milestone
		err := rows.Scan(
			&milestone.ID, &milestone.GoalID, &milestone.Title, &milestone.Description,
			&milestone.TargetDate, &milestone.Completed, &milestone.CompletedAt,
			&milestone.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}
		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

func (r *milestoneRepository) GetUpcomingByUserID(ctx context.Context, userID entities.UserID, before time.Time) ([]*entities.Milestone, error) {
	query := `
		SELECT m.id, m.goal_id, m.title, m.description, m.target_date,
			   m.completed, m.completed_at, m.created_at
		FROM milestones m
		JOIN goals g ON m.goal_id = g.id
		WHERE g.user_id = $1 AND m.target_date <= $2 AND m.completed = false
		ORDER BY m.target_date ASC`

	rows, err := r.pool.Query(ctx, query, userID, before)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming milestones: %w", err)
	}
	defer rows.Close()

	var milestones []*entities.Milestone
	for rows.Next() {
		var milestone entities.Milestone
		err := rows.Scan(
			&milestone.ID, &milestone.GoalID, &milestone.Title, &milestone.Description,
			&milestone.TargetDate, &milestone.Completed, &milestone.CompletedAt,
			&milestone.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}
		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

func (r *milestoneRepository) Update(ctx context.Context, milestone *entities.Milestone) error {
	query := `
		UPDATE milestones 
		SET title = $2, description = $3, target_date = $4,
			completed = $5, completed_at = $6
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query,
		milestone.ID, milestone.Title, milestone.Description,
		milestone.TargetDate, milestone.Completed, milestone.CompletedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}

	return nil
}

func (r *milestoneRepository) Delete(ctx context.Context, id entities.MilestoneID) error {
	query := `DELETE FROM milestones WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}

	return nil
}

func (r *milestoneRepository) MarkCompleted(ctx context.Context, milestoneID entities.MilestoneID) error {
	query := `
		UPDATE milestones 
		SET completed = true, completed_at = $2
		WHERE id = $1`

	now := time.Now()
	_, err := r.pool.Exec(ctx, query, milestoneID, now)
	if err != nil {
		return fmt.Errorf("failed to mark milestone as completed: %w", err)
	}

	return nil
}