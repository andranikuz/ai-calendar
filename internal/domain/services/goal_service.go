package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type GoalService struct{}

func NewGoalService() *GoalService {
	return &GoalService{}
}

// ValidateGoalCreation validates goal data before creation
func (s *GoalService) ValidateGoalCreation(goal *entities.Goal) error {
	if err := s.ValidateGoalTitle(goal.Title); err != nil {
		return fmt.Errorf("invalid title: %w", err)
	}
	
	if err := s.ValidateGoalDescription(goal.Description); err != nil {
		return fmt.Errorf("invalid description: %w", err)
	}
	
	if err := s.ValidateGoalCategory(goal.Category); err != nil {
		return fmt.Errorf("invalid category: %w", err)
	}
	
	if err := s.ValidateGoalPriority(goal.Priority); err != nil {
		return fmt.Errorf("invalid priority: %w", err)
	}
	
	if goal.Deadline != nil {
		if err := s.ValidateGoalDeadline(*goal.Deadline); err != nil {
			return fmt.Errorf("invalid deadline: %w", err)
		}
	}
	
	if err := s.ValidateGoalProgress(goal.Progress); err != nil {
		return fmt.Errorf("invalid progress: %w", err)
	}
	
	return nil
}

// ValidateGoalTitle validates goal title
func (s *GoalService) ValidateGoalTitle(title string) error {
	if title == "" {
		return fmt.Errorf("title is required")
	}
	
	if len(title) < 2 {
		return fmt.Errorf("title is too short (min 2 characters)")
	}
	
	if len(title) > 255 {
		return fmt.Errorf("title is too long (max 255 characters)")
	}
	
	return nil
}

// ValidateGoalDescription validates goal description
func (s *GoalService) ValidateGoalDescription(description string) error {
	if len(description) > 1000 {
		return fmt.Errorf("description is too long (max 1000 characters)")
	}
	
	return nil
}

// ValidateGoalCategory validates goal category
func (s *GoalService) ValidateGoalCategory(category entities.GoalCategory) error {
	validCategories := []entities.GoalCategory{
		entities.GoalCategoryHealth,
		entities.GoalCategoryCareer,
		entities.GoalCategoryEducation,
		entities.GoalCategoryPersonal,
		entities.GoalCategoryFinancial,
		entities.GoalCategoryRelationship,
	}
	
	for _, validCategory := range validCategories {
		if category == validCategory {
			return nil
		}
	}
	
	return fmt.Errorf("invalid category: %s", category)
}

// ValidateGoalPriority validates goal priority
func (s *GoalService) ValidateGoalPriority(priority entities.Priority) error {
	if !priority.IsValid() {
		return fmt.Errorf("invalid priority: %d", priority)
	}
	
	return nil
}

// ValidateGoalDeadline validates goal deadline
func (s *GoalService) ValidateGoalDeadline(deadline time.Time) error {
	now := time.Now()
	
	if deadline.Before(now) {
		return fmt.Errorf("deadline cannot be in the past")
	}
	
	// Check if deadline is too far in the future (e.g., 10 years)
	maxDeadline := now.AddDate(10, 0, 0)
	if deadline.After(maxDeadline) {
		return fmt.Errorf("deadline cannot be more than 10 years in the future")
	}
	
	return nil
}

// ValidateGoalProgress validates goal progress
func (s *GoalService) ValidateGoalProgress(progress int) error {
	if progress < 0 {
		return fmt.Errorf("progress cannot be negative")
	}
	
	if progress > 100 {
		return fmt.Errorf("progress cannot exceed 100%%")
	}
	
	return nil
}

// ValidateTaskCreation validates task data before creation
func (s *GoalService) ValidateTaskCreation(task *entities.Task) error {
	if err := s.ValidateTaskTitle(task.Title); err != nil {
		return fmt.Errorf("invalid title: %w", err)
	}
	
	if err := s.ValidateTaskDescription(task.Description); err != nil {
		return fmt.Errorf("invalid description: %w", err)
	}
	
	if err := s.ValidateGoalPriority(task.Priority); err != nil {
		return fmt.Errorf("invalid priority: %w", err)
	}
	
	if err := s.ValidateTaskDuration(task.EstimatedDuration); err != nil {
		return fmt.Errorf("invalid estimated duration: %w", err)
	}
	
	if task.DueDate != nil {
		if err := s.ValidateTaskDueDate(*task.DueDate); err != nil {
			return fmt.Errorf("invalid due date: %w", err)
		}
	}
	
	return nil
}

// ValidateTaskTitle validates task title
func (s *GoalService) ValidateTaskTitle(title string) error {
	return s.ValidateGoalTitle(title) // Same rules as goal title
}

// ValidateTaskDescription validates task description
func (s *GoalService) ValidateTaskDescription(description string) error {
	return s.ValidateGoalDescription(description) // Same rules as goal description
}

// ValidateTaskDuration validates task estimated duration
func (s *GoalService) ValidateTaskDuration(duration int) error {
	if duration <= 0 {
		return fmt.Errorf("estimated duration must be positive")
	}
	
	// Maximum 8 hours (480 minutes)
	if duration > 480 {
		return fmt.Errorf("estimated duration cannot exceed 8 hours")
	}
	
	return nil
}

// ValidateTaskDueDate validates task due date
func (s *GoalService) ValidateTaskDueDate(dueDate time.Time) error {
	return s.ValidateGoalDeadline(dueDate) // Same rules as goal deadline
}

// ValidateMilestoneCreation validates milestone data
func (s *GoalService) ValidateMilestoneCreation(milestone *entities.Milestone) error {
	if err := s.ValidateMilestoneTitle(milestone.Title); err != nil {
		return fmt.Errorf("invalid title: %w", err)
	}
	
	if err := s.ValidateMilestoneDescription(milestone.Description); err != nil {
		return fmt.Errorf("invalid description: %w", err)
	}
	
	if err := s.ValidateMilestoneTargetDate(milestone.TargetDate); err != nil {
		return fmt.Errorf("invalid target date: %w", err)
	}
	
	return nil
}

// ValidateMilestoneTitle validates milestone title
func (s *GoalService) ValidateMilestoneTitle(title string) error {
	return s.ValidateGoalTitle(title) // Same rules as goal title
}

// ValidateMilestoneDescription validates milestone description
func (s *GoalService) ValidateMilestoneDescription(description string) error {
	return s.ValidateGoalDescription(description) // Same rules as goal description
}

// ValidateMilestoneTargetDate validates milestone target date
func (s *GoalService) ValidateMilestoneTargetDate(targetDate time.Time) error {
	return s.ValidateGoalDeadline(targetDate) // Same rules as goal deadline
}

// CalculateGoalProgress calculates goal progress based on completed tasks and milestones
func (s *GoalService) CalculateGoalProgress(tasks []*entities.Task, milestones []*entities.Milestone) int {
	if len(tasks) == 0 && len(milestones) == 0 {
		return 0
	}
	
	totalItems := len(tasks) + len(milestones)
	completedItems := 0
	
	// Count completed tasks
	for _, task := range tasks {
		if task.Status == entities.TaskStatusCompleted {
			completedItems++
		}
	}
	
	// Count completed milestones
	for _, milestone := range milestones {
		if milestone.Completed {
			completedItems++
		}
	}
	
	return (completedItems * 100) / totalItems
}

// SanitizeGoalTitle sanitizes goal title
func (s *GoalService) SanitizeGoalTitle(title string) string {
	return strings.TrimSpace(title)
}

// SanitizeGoalDescription sanitizes goal description
func (s *GoalService) SanitizeGoalDescription(description string) string {
	return strings.TrimSpace(description)
}

// IsGoalOverdue checks if goal is overdue
func (s *GoalService) IsGoalOverdue(goal *entities.Goal) bool {
	if goal.Deadline == nil {
		return false
	}
	
	return time.Now().After(*goal.Deadline) && goal.Status != entities.GoalStatusCompleted
}

// IsTaskOverdue checks if task is overdue
func (s *GoalService) IsTaskOverdue(task *entities.Task) bool {
	if task.DueDate == nil {
		return false
	}
	
	return time.Now().After(*task.DueDate) && task.Status != entities.TaskStatusCompleted
}

// IsMilestoneOverdue checks if milestone is overdue
func (s *GoalService) IsMilestoneOverdue(milestone *entities.Milestone) bool {
	return time.Now().After(milestone.TargetDate) && !milestone.Completed
}

// CanCompleteGoal checks if goal can be marked as completed
func (s *GoalService) CanCompleteGoal(ctx context.Context, goal *entities.Goal, tasks []*entities.Task, milestones []*entities.Milestone) error {
	// Check if all required tasks are completed
	for _, task := range tasks {
		if task.Status != entities.TaskStatusCompleted && task.Status != entities.TaskStatusCancelled {
			return fmt.Errorf("cannot complete goal: task '%s' is not completed", task.Title)
		}
	}
	
	// Check if all required milestones are completed
	for _, milestone := range milestones {
		if !milestone.Completed {
			return fmt.Errorf("cannot complete goal: milestone '%s' is not completed", milestone.Title)
		}
	}
	
	return nil
}

// SuggestGoalCategory suggests category based on goal title and description
func (s *GoalService) SuggestGoalCategory(title, description string) entities.GoalCategory {
	text := strings.ToLower(title + " " + description)
	
	// Health-related keywords
	healthKeywords := []string{"health", "fitness", "exercise", "diet", "weight", "medical", "doctor", "gym", "workout", "nutrition"}
	for _, keyword := range healthKeywords {
		if strings.Contains(text, keyword) {
			return entities.GoalCategoryHealth
		}
	}
	
	// Career-related keywords
	careerKeywords := []string{"career", "job", "work", "promotion", "salary", "skill", "training", "certification", "professional", "business"}
	for _, keyword := range careerKeywords {
		if strings.Contains(text, keyword) {
			return entities.GoalCategoryCareer
		}
	}
	
	// Education-related keywords
	educationKeywords := []string{"education", "study", "learn", "course", "degree", "school", "university", "book", "knowledge", "academic"}
	for _, keyword := range educationKeywords {
		if strings.Contains(text, keyword) {
			return entities.GoalCategoryEducation
		}
	}
	
	// Financial-related keywords
	financialKeywords := []string{"money", "financial", "save", "invest", "budget", "debt", "income", "expense", "retirement", "fund"}
	for _, keyword := range financialKeywords {
		if strings.Contains(text, keyword) {
			return entities.GoalCategoryFinancial
		}
	}
	
	// Relationship-related keywords
	relationshipKeywords := []string{"relationship", "family", "friend", "social", "love", "marriage", "dating", "communication", "network"}
	for _, keyword := range relationshipKeywords {
		if strings.Contains(text, keyword) {
			return entities.GoalCategoryRelationship
		}
	}
	
	// Default to personal
	return entities.GoalCategoryPersonal
}