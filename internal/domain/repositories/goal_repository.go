package repositories

import (
	"context"
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type GoalRepository interface {
	// Create a new goal
	Create(ctx context.Context, goal *entities.Goal) error
	
	// Get goal by ID
	GetByID(ctx context.Context, id entities.GoalID) (*entities.Goal, error)
	
	// Get all goals for a user
	GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error)
	
	// Get goals by status for a user
	GetByUserIDAndStatus(ctx context.Context, userID entities.UserID, status entities.GoalStatus) ([]*entities.Goal, error)
	
	// Get goals by category for a user
	GetByUserIDAndCategory(ctx context.Context, userID entities.UserID, category entities.GoalCategory) ([]*entities.Goal, error)
	
	// Get goals with upcoming deadlines
	GetByDeadlineBefore(ctx context.Context, userID entities.UserID, deadline time.Time) ([]*entities.Goal, error)
	
	// Update goal
	Update(ctx context.Context, goal *entities.Goal) error
	
	// Delete goal
	Delete(ctx context.Context, id entities.GoalID) error
	
	// Update goal progress
	UpdateProgress(ctx context.Context, goalID entities.GoalID, progress int) error
	
	// Get goals with pagination
	GetByUserIDPaginated(ctx context.Context, userID entities.UserID, offset, limit int) ([]*entities.Goal, int64, error)
}

type TaskRepository interface {
	// Create a new task
	Create(ctx context.Context, task *entities.Task) error
	
	// Get task by ID
	GetByID(ctx context.Context, id entities.TaskID) (*entities.Task, error)
	
	// Get all tasks for a goal
	GetByGoalID(ctx context.Context, goalID entities.GoalID) ([]*entities.Task, error)
	
	// Get tasks by status for a goal
	GetByGoalIDAndStatus(ctx context.Context, goalID entities.GoalID, status entities.TaskStatus) ([]*entities.Task, error)
	
	// Get all tasks for a user (across all goals)
	GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Task, error)
	
	// Get tasks with upcoming due dates
	GetByDueDateBefore(ctx context.Context, userID entities.UserID, dueDate time.Time) ([]*entities.Task, error)
	
	// Update task
	Update(ctx context.Context, task *entities.Task) error
	
	// Delete task
	Delete(ctx context.Context, id entities.TaskID) error
	
	// Mark task as completed
	MarkCompleted(ctx context.Context, taskID entities.TaskID) error
	
	// Get tasks with pagination
	GetByUserIDPaginated(ctx context.Context, userID entities.UserID, offset, limit int) ([]*entities.Task, int64, error)
}

type MilestoneRepository interface {
	// Create a new milestone
	Create(ctx context.Context, milestone *entities.Milestone) error
	
	// Get milestone by ID
	GetByID(ctx context.Context, id entities.MilestoneID) (*entities.Milestone, error)
	
	// Get all milestones for a goal
	GetByGoalID(ctx context.Context, goalID entities.GoalID) ([]*entities.Milestone, error)
	
	// Get completed milestones for a goal
	GetCompletedByGoalID(ctx context.Context, goalID entities.GoalID) ([]*entities.Milestone, error)
	
	// Get upcoming milestones for a user
	GetUpcomingByUserID(ctx context.Context, userID entities.UserID, before time.Time) ([]*entities.Milestone, error)
	
	// Update milestone
	Update(ctx context.Context, milestone *entities.Milestone) error
	
	// Delete milestone
	Delete(ctx context.Context, id entities.MilestoneID) error
	
	// Mark milestone as completed
	MarkCompleted(ctx context.Context, milestoneID entities.MilestoneID) error
}