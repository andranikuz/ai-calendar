package commands

import (
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

// CreateGoalCommand represents a command to create a new goal
type CreateGoalCommand struct {
	UserID      entities.UserID        `json:"user_id" validate:"required"`
	Title       string                 `json:"title" validate:"required,min=2,max=255"`
	Description string                 `json:"description" validate:"max=1000"`
	Category    entities.GoalCategory  `json:"category" validate:"required"`
	Priority    entities.Priority      `json:"priority" validate:"required"`
	Deadline    *time.Time             `json:"deadline,omitempty"`
}

// UpdateGoalCommand represents a command to update a goal
type UpdateGoalCommand struct {
	GoalID      entities.GoalID        `json:"goal_id" validate:"required"`
	UserID      entities.UserID        `json:"user_id" validate:"required"`
	Title       *string                `json:"title,omitempty" validate:"omitempty,min=2,max=255"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Category    *entities.GoalCategory `json:"category,omitempty"`
	Priority    *entities.Priority     `json:"priority,omitempty"`
	Status      *entities.GoalStatus   `json:"status,omitempty"`
	Progress    *int                   `json:"progress,omitempty" validate:"omitempty,min=0,max=100"`
	Deadline    *time.Time             `json:"deadline,omitempty"`
}

// DeleteGoalCommand represents a command to delete a goal
type DeleteGoalCommand struct {
	GoalID entities.GoalID `json:"goal_id" validate:"required"`
	UserID entities.UserID `json:"user_id" validate:"required"`
}

// UpdateGoalProgressCommand represents a command to update goal progress
type UpdateGoalProgressCommand struct {
	GoalID   entities.GoalID `json:"goal_id" validate:"required"`
	UserID   entities.UserID `json:"user_id" validate:"required"`
	Progress int             `json:"progress" validate:"min=0,max=100"`
}

// CreateTaskCommand represents a command to create a task for a goal
type CreateTaskCommand struct {
	GoalID            entities.GoalID    `json:"goal_id" validate:"required"`
	UserID            entities.UserID    `json:"user_id" validate:"required"`
	Title             string             `json:"title" validate:"required,min=2,max=255"`
	Description       string             `json:"description" validate:"max=1000"`
	Priority          entities.Priority  `json:"priority" validate:"required"`
	EstimatedDuration int                `json:"estimated_duration" validate:"min=1"` // minutes
	DueDate           *time.Time         `json:"due_date,omitempty"`
}

// UpdateTaskCommand represents a command to update a task
type UpdateTaskCommand struct {
	TaskID            entities.TaskID    `json:"task_id" validate:"required"`
	UserID            entities.UserID    `json:"user_id" validate:"required"`
	Title             *string            `json:"title,omitempty" validate:"omitempty,min=2,max=255"`
	Description       *string            `json:"description,omitempty" validate:"omitempty,max=1000"`
	Priority          *entities.Priority `json:"priority,omitempty"`
	Status            *entities.TaskStatus `json:"status,omitempty"`
	EstimatedDuration *int               `json:"estimated_duration,omitempty" validate:"omitempty,min=1"`
	DueDate           *time.Time         `json:"due_date,omitempty"`
}

// CompleteTaskCommand represents a command to mark a task as completed
type CompleteTaskCommand struct {
	TaskID entities.TaskID `json:"task_id" validate:"required"`
	UserID entities.UserID `json:"user_id" validate:"required"`
}

// DeleteTaskCommand represents a command to delete a task
type DeleteTaskCommand struct {
	TaskID entities.TaskID `json:"task_id" validate:"required"`
	UserID entities.UserID `json:"user_id" validate:"required"`
}

// CreateMilestoneCommand represents a command to create a milestone
type CreateMilestoneCommand struct {
	GoalID      entities.GoalID `json:"goal_id" validate:"required"`
	UserID      entities.UserID `json:"user_id" validate:"required"`
	Title       string          `json:"title" validate:"required,min=2,max=255"`
	Description string          `json:"description" validate:"max=1000"`
	TargetDate  time.Time       `json:"target_date" validate:"required"`
}

// CompleteMilestoneCommand represents a command to mark a milestone as completed
type CompleteMilestoneCommand struct {
	MilestoneID entities.MilestoneID `json:"milestone_id" validate:"required"`
	UserID      entities.UserID      `json:"user_id" validate:"required"`
}

// Results
type CreateGoalResult struct {
	GoalID    entities.GoalID `json:"goal_id"`
	CreatedAt time.Time       `json:"created_at"`
}

type UpdateGoalResult struct {
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaskResult struct {
	TaskID    entities.TaskID `json:"task_id"`
	CreatedAt time.Time       `json:"created_at"`
}

type CreateMilestoneResult struct {
	MilestoneID entities.MilestoneID `json:"milestone_id"`
	CreatedAt   time.Time            `json:"created_at"`
}