package entities

import (
	"time"
)

type GoalID string

type Goal struct {
	ID          GoalID    `json:"id"`
	UserID      UserID    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    GoalCategory `json:"category"`
	Priority    Priority  `json:"priority"`
	Status      GoalStatus `json:"status"`
	Progress    int       `json:"progress"` // 0-100
	Deadline    *time.Time `json:"deadline"`
	Milestones  []Milestone `json:"milestones"`
	Tasks       []Task    `json:"tasks"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GoalCategory string

const (
	GoalCategoryHealth     GoalCategory = "health"
	GoalCategoryCareer     GoalCategory = "career"
	GoalCategoryEducation  GoalCategory = "education"
	GoalCategoryPersonal   GoalCategory = "personal"
	GoalCategoryFinancial  GoalCategory = "financial"
	GoalCategoryRelationship GoalCategory = "relationship"
)

type GoalStatus string

const (
	GoalStatusDraft      GoalStatus = "draft"
	GoalStatusActive     GoalStatus = "active"
	GoalStatusPaused     GoalStatus = "paused"
	GoalStatusCompleted  GoalStatus = "completed"
	GoalStatusCancelled  GoalStatus = "cancelled"
)

type MilestoneID string

type Milestone struct {
	ID          MilestoneID `json:"id"`
	GoalID      GoalID     `json:"goal_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	TargetDate  time.Time  `json:"target_date"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type TaskID string

type Task struct {
	ID          TaskID    `json:"id"`
	GoalID      GoalID    `json:"goal_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	Status      TaskStatus `json:"status"`
	EstimatedDuration int `json:"estimated_duration"` // minutes
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)