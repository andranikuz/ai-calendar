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

type taskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) repositories.TaskRepository {
	return &taskRepository{pool: pool}
}

func (r *taskRepository) Create(ctx context.Context, task *entities.Task) error {
	query := `
		INSERT INTO tasks (
			id, goal_id, title, description, priority, status,
			estimated_duration, due_date, completed_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.pool.Exec(ctx, query,
		task.ID, task.GoalID, task.Title, task.Description,
		task.Priority, task.Status, task.EstimatedDuration,
		task.DueDate, task.CompletedAt, task.CreatedAt, task.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	
	return nil
}

func (r *taskRepository) GetByID(ctx context.Context, id entities.TaskID) (*entities.Task, error) {
	query := `
		SELECT id, goal_id, title, description, priority, status,
			   estimated_duration, due_date, completed_at, created_at, updated_at
		FROM tasks 
		WHERE id = $1`

	var task entities.Task
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&task.ID, &task.GoalID, &task.Title, &task.Description,
		&task.Priority, &task.Status, &task.EstimatedDuration,
		&task.DueDate, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task by ID: %w", err)
	}

	return &task, nil
}

func (r *taskRepository) GetByGoalID(ctx context.Context, goalID entities.GoalID) ([]*entities.Task, error) {
	query := `
		SELECT id, goal_id, title, description, priority, status,
			   estimated_duration, due_date, completed_at, created_at, updated_at
		FROM tasks 
		WHERE goal_id = $1
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by goal ID: %w", err)
	}
	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(
			&task.ID, &task.GoalID, &task.Title, &task.Description,
			&task.Priority, &task.Status, &task.EstimatedDuration,
			&task.DueDate, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *taskRepository) GetByGoalIDAndStatus(ctx context.Context, goalID entities.GoalID, status entities.TaskStatus) ([]*entities.Task, error) {
	query := `
		SELECT id, goal_id, title, description, priority, status,
			   estimated_duration, due_date, completed_at, created_at, updated_at
		FROM tasks 
		WHERE goal_id = $1 AND status = $2
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, goalID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}
	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(
			&task.ID, &task.GoalID, &task.Title, &task.Description,
			&task.Priority, &task.Status, &task.EstimatedDuration,
			&task.DueDate, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *taskRepository) GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Task, error) {
	query := `
		SELECT t.id, t.goal_id, t.title, t.description, t.priority, t.status,
			   t.estimated_duration, t.due_date, t.completed_at, t.created_at, t.updated_at
		FROM tasks t
		JOIN goals g ON t.goal_id = g.id
		WHERE g.user_id = $1
		ORDER BY t.created_at DESC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by user ID: %w", err)
	}
	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(
			&task.ID, &task.GoalID, &task.Title, &task.Description,
			&task.Priority, &task.Status, &task.EstimatedDuration,
			&task.DueDate, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *taskRepository) GetByDueDateBefore(ctx context.Context, userID entities.UserID, dueDate time.Time) ([]*entities.Task, error) {
	query := `
		SELECT t.id, t.goal_id, t.title, t.description, t.priority, t.status,
			   t.estimated_duration, t.due_date, t.completed_at, t.created_at, t.updated_at
		FROM tasks t
		JOIN goals g ON t.goal_id = g.id
		WHERE g.user_id = $1 AND t.due_date <= $2 AND t.status != 'completed'
		ORDER BY t.due_date ASC`

	rows, err := r.pool.Query(ctx, query, userID, dueDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by due date: %w", err)
	}
	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(
			&task.ID, &task.GoalID, &task.Title, &task.Description,
			&task.Priority, &task.Status, &task.EstimatedDuration,
			&task.DueDate, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *entities.Task) error {
	query := `
		UPDATE tasks 
		SET title = $2, description = $3, priority = $4, status = $5,
			estimated_duration = $6, due_date = $7, completed_at = $8, updated_at = $9
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query,
		task.ID, task.Title, task.Description, task.Priority,
		task.Status, task.EstimatedDuration, task.DueDate,
		task.CompletedAt, task.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id entities.TaskID) error {
	query := `DELETE FROM tasks WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (r *taskRepository) MarkCompleted(ctx context.Context, taskID entities.TaskID) error {
	query := `
		UPDATE tasks 
		SET status = 'completed', completed_at = $2, updated_at = $2
		WHERE id = $1`

	now := time.Now()
	_, err := r.pool.Exec(ctx, query, taskID, now)
	if err != nil {
		return fmt.Errorf("failed to mark task as completed: %w", err)
	}

	return nil
}

func (r *taskRepository) GetByUserIDPaginated(ctx context.Context, userID entities.UserID, offset, limit int) ([]*entities.Task, int64, error) {
	// Get total count
	countQuery := `
		SELECT COUNT(*) 
		FROM tasks t
		JOIN goals g ON t.goal_id = g.id
		WHERE g.user_id = $1`
	var totalCount int64
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tasks count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT t.id, t.goal_id, t.title, t.description, t.priority, t.status,
			   t.estimated_duration, t.due_date, t.completed_at, t.created_at, t.updated_at
		FROM tasks t
		JOIN goals g ON t.goal_id = g.id
		WHERE g.user_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(
			&task.ID, &task.GoalID, &task.Title, &task.Description,
			&task.Priority, &task.Status, &task.EstimatedDuration,
			&task.DueDate, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, totalCount, nil
}