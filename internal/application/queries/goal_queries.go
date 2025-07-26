package queries

import (
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

// GetGoalByIDQuery represents a query to get goal by ID
type GetGoalByIDQuery struct {
	GoalID entities.GoalID `json:"goal_id" validate:"required"`
	UserID entities.UserID `json:"user_id" validate:"required"`
}

// GetGoalsByUserIDQuery represents a query to get all goals for a user
type GetGoalsByUserIDQuery struct {
	UserID entities.UserID `json:"user_id" validate:"required"`
	Offset int             `json:"offset" validate:"min=0"`
	Limit  int             `json:"limit" validate:"min=1,max=100"`
}

// GetGoalsByStatusQuery represents a query to get goals by status
type GetGoalsByStatusQuery struct {
	UserID entities.UserID    `json:"user_id" validate:"required"`
	Status entities.GoalStatus `json:"status" validate:"required"`
	Offset int                 `json:"offset" validate:"min=0"`
	Limit  int                 `json:"limit" validate:"min=1,max=100"`
}

// GetGoalsByCategoryQuery represents a query to get goals by category
type GetGoalsByCategoryQuery struct {
	UserID   entities.UserID       `json:"user_id" validate:"required"`
	Category entities.GoalCategory `json:"category" validate:"required"`
	Offset   int                   `json:"offset" validate:"min=0"`
	Limit    int                   `json:"limit" validate:"min=1,max=100"`
}

// GetGoalsWithUpcomingDeadlinesQuery represents a query to get goals with upcoming deadlines
type GetGoalsWithUpcomingDeadlinesQuery struct {
	UserID   entities.UserID `json:"user_id" validate:"required"`
	Deadline time.Time       `json:"deadline" validate:"required"`
	Limit    int             `json:"limit" validate:"min=1,max=100"`
}

// GetTasksByGoalIDQuery represents a query to get tasks for a goal
type GetTasksByGoalIDQuery struct {
	GoalID entities.GoalID `json:"goal_id" validate:"required"`
	UserID entities.UserID `json:"user_id" validate:"required"`
}

// GetTasksByUserIDQuery represents a query to get all tasks for a user
type GetTasksByUserIDQuery struct {
	UserID entities.UserID `json:"user_id" validate:"required"`
	Offset int             `json:"offset" validate:"min=0"`
	Limit  int             `json:"limit" validate:"min=1,max=100"`
}

// GetTasksByStatusQuery represents a query to get tasks by status
type GetTasksByStatusQuery struct {
	UserID entities.UserID     `json:"user_id" validate:"required"`
	Status entities.TaskStatus `json:"status" validate:"required"`
	Offset int                 `json:"offset" validate:"min=0"`
	Limit  int                 `json:"limit" validate:"min=1,max=100"`
}

// GetUpcomingTasksQuery represents a query to get upcoming tasks
type GetUpcomingTasksQuery struct {
	UserID  entities.UserID `json:"user_id" validate:"required"`
	DueDate time.Time       `json:"due_date" validate:"required"`
	Limit   int             `json:"limit" validate:"min=1,max=100"`
}

// GetMilestonesByGoalIDQuery represents a query to get milestones for a goal
type GetMilestonesByGoalIDQuery struct {
	GoalID entities.GoalID `json:"goal_id" validate:"required"`
	UserID entities.UserID `json:"user_id" validate:"required"`
}

// GetUpcomingMilestonesQuery represents a query to get upcoming milestones
type GetUpcomingMilestonesQuery struct {
	UserID     entities.UserID `json:"user_id" validate:"required"`
	TargetDate time.Time       `json:"target_date" validate:"required"`
	Limit      int             `json:"limit" validate:"min=1,max=100"`
}

// Results
type GetGoalResult struct {
	Goal *entities.Goal `json:"goal"`
}

type GetGoalsResult struct {
	Goals      []*entities.Goal `json:"goals"`
	TotalCount int64            `json:"total_count"`
	Offset     int              `json:"offset"`
	Limit      int              `json:"limit"`
}

type GetTasksResult struct {
	Tasks      []*entities.Task `json:"tasks"`
	TotalCount int64            `json:"total_count"`
	Offset     int              `json:"offset"`
	Limit      int              `json:"limit"`
}

type GetMilestonesResult struct {
	Milestones []*entities.Milestone `json:"milestones"`
}

type GetGoalStatsResult struct {
	TotalGoals      int                                  `json:"total_goals"`
	CompletedGoals  int                                  `json:"completed_goals"`
	ActiveGoals     int                                  `json:"active_goals"`
	GoalsByCategory map[entities.GoalCategory]int        `json:"goals_by_category"`
	GoalsByPriority map[entities.Priority]int            `json:"goals_by_priority"`
	OverdueGoals    int                                  `json:"overdue_goals"`
}