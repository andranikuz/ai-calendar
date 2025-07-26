package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/andranikuz/smart-goal-calendar/internal/application/commands"
	"github.com/andranikuz/smart-goal-calendar/internal/application/queries"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/services"
)

type GoalHandler struct {
	goalRepo      repositories.GoalRepository
	taskRepo      repositories.TaskRepository
	milestoneRepo repositories.MilestoneRepository
	goalService   *services.GoalService
}

func NewGoalHandler(
	goalRepo repositories.GoalRepository,
	taskRepo repositories.TaskRepository,
	milestoneRepo repositories.MilestoneRepository,
	goalService *services.GoalService,
) *GoalHandler {
	return &GoalHandler{
		goalRepo:      goalRepo,
		taskRepo:      taskRepo,
		milestoneRepo: milestoneRepo,
		goalService:   goalService,
	}
}

// Command Handlers

func (h *GoalHandler) HandleCreateGoal(ctx context.Context, cmd commands.CreateGoalCommand) (*commands.CreateGoalResult, error) {
	// Create goal entity
	now := time.Now()
	goal := &entities.Goal{
		ID:          entities.GoalID(uuid.New().String()),
		UserID:      cmd.UserID,
		Title:       h.goalService.SanitizeGoalTitle(cmd.Title),
		Description: h.goalService.SanitizeGoalDescription(cmd.Description),
		Category:    cmd.Category,
		Priority:    cmd.Priority,
		Status:      entities.GoalStatusDraft,
		Progress:    0,
		Deadline:    cmd.Deadline,
		Milestones:  []entities.Milestone{},
		Tasks:       []entities.Task{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	// Validate goal
	if err := h.goalService.ValidateGoalCreation(goal); err != nil {
		return nil, fmt.Errorf("goal validation failed: %w", err)
	}
	
	// Save goal
	if err := h.goalRepo.Create(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}
	
	return &commands.CreateGoalResult{
		GoalID:    goal.ID,
		CreatedAt: goal.CreatedAt,
	}, nil
}

func (h *GoalHandler) HandleUpdateGoal(ctx context.Context, cmd commands.UpdateGoalCommand) (*commands.UpdateGoalResult, error) {
	// Get existing goal
	goal, err := h.goalRepo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}
	
	// Check ownership
	if goal.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: goal belongs to different user")
	}
	
	// Update fields if provided
	if cmd.Title != nil {
		goal.Title = h.goalService.SanitizeGoalTitle(*cmd.Title)
	}
	
	if cmd.Description != nil {
		goal.Description = h.goalService.SanitizeGoalDescription(*cmd.Description)
	}
	
	if cmd.Category != nil {
		goal.Category = *cmd.Category
	}
	
	if cmd.Priority != nil {
		goal.Priority = *cmd.Priority
	}
	
	if cmd.Status != nil {
		goal.Status = *cmd.Status
	}
	
	if cmd.Progress != nil {
		goal.Progress = *cmd.Progress
	}
	
	if cmd.Deadline != nil {
		goal.Deadline = cmd.Deadline
	}
	
	// Update timestamp
	goal.UpdatedAt = time.Now()
	
	// Validate updated goal
	if err := h.goalService.ValidateGoalCreation(goal); err != nil {
		return nil, fmt.Errorf("goal validation failed: %w", err)
	}
	
	// Save updated goal
	if err := h.goalRepo.Update(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}
	
	return &commands.UpdateGoalResult{
		UpdatedAt: goal.UpdatedAt,
	}, nil
}

func (h *GoalHandler) HandleDeleteGoal(ctx context.Context, cmd commands.DeleteGoalCommand) error {
	// Get goal to check ownership
	goal, err := h.goalRepo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal == nil {
		return fmt.Errorf("goal not found")
	}
	
	// Check ownership
	if goal.UserID != cmd.UserID {
		return fmt.Errorf("access denied: goal belongs to different user")
	}
	
	// Delete goal (cascading deletes will handle tasks and milestones)
	if err := h.goalRepo.Delete(ctx, cmd.GoalID); err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}
	
	return nil
}

func (h *GoalHandler) HandleUpdateGoalProgress(ctx context.Context, cmd commands.UpdateGoalProgressCommand) (*commands.UpdateGoalResult, error) {
	// Get existing goal
	goal, err := h.goalRepo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}
	
	// Check ownership
	if goal.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: goal belongs to different user")
	}
	
	// Update progress
	if err := h.goalRepo.UpdateProgress(ctx, cmd.GoalID, cmd.Progress); err != nil {
		return nil, fmt.Errorf("failed to update goal progress: %w", err)
	}
	
	// Auto-complete goal if progress is 100%
	if cmd.Progress == 100 && goal.Status != entities.GoalStatusCompleted {
		goal.Status = entities.GoalStatusCompleted
		goal.UpdatedAt = time.Now()
		if err := h.goalRepo.Update(ctx, goal); err != nil {
			return nil, fmt.Errorf("failed to mark goal as completed: %w", err)
		}
	}
	
	return &commands.UpdateGoalResult{
		UpdatedAt: time.Now(),
	}, nil
}

func (h *GoalHandler) HandleCreateTask(ctx context.Context, cmd commands.CreateTaskCommand) (*commands.CreateTaskResult, error) {
	// Verify goal exists and user has access
	goal, err := h.goalRepo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}
	
	if goal.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: goal belongs to different user")
	}
	
	// Create task entity
	now := time.Now()
	task := &entities.Task{
		ID:                entities.TaskID(uuid.New().String()),
		GoalID:            cmd.GoalID,
		Title:             h.goalService.SanitizeGoalTitle(cmd.Title),
		Description:       h.goalService.SanitizeGoalDescription(cmd.Description),
		Priority:          cmd.Priority,
		Status:            entities.TaskStatusPending,
		EstimatedDuration: cmd.EstimatedDuration,
		DueDate:           cmd.DueDate,
		CompletedAt:       nil,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	
	// Validate task
	if err := h.goalService.ValidateTaskCreation(task); err != nil {
		return nil, fmt.Errorf("task validation failed: %w", err)
	}
	
	// Save task
	if err := h.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	
	return &commands.CreateTaskResult{
		TaskID:    task.ID,
		CreatedAt: task.CreatedAt,
	}, nil
}

func (h *GoalHandler) HandleCompleteTask(ctx context.Context, cmd commands.CompleteTaskCommand) error {
	// Get task to check ownership
	task, err := h.taskRepo.GetByID(ctx, cmd.TaskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	
	if task == nil {
		return fmt.Errorf("task not found")
	}
	
	// Get goal to check ownership
	goal, err := h.goalRepo.GetByID(ctx, task.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal.UserID != cmd.UserID {
		return fmt.Errorf("access denied: task belongs to different user")
	}
	
	// Mark task as completed
	if err := h.taskRepo.MarkCompleted(ctx, cmd.TaskID); err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}
	
	// Recalculate goal progress
	tasks, err := h.taskRepo.GetByGoalID(ctx, task.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal tasks: %w", err)
	}
	
	milestones, err := h.milestoneRepo.GetByGoalID(ctx, task.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal milestones: %w", err)
	}
	
	newProgress := h.goalService.CalculateGoalProgress(tasks, milestones)
	if err := h.goalRepo.UpdateProgress(ctx, task.GoalID, newProgress); err != nil {
		return fmt.Errorf("failed to update goal progress: %w", err)
	}
	
	return nil
}

func (h *GoalHandler) HandleCreateMilestone(ctx context.Context, cmd commands.CreateMilestoneCommand) (*commands.CreateMilestoneResult, error) {
	// Verify goal exists and user has access
	goal, err := h.goalRepo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}
	
	if goal.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: goal belongs to different user")
	}
	
	// Create milestone entity
	now := time.Now()
	milestone := &entities.Milestone{
		ID:          entities.MilestoneID(uuid.New().String()),
		GoalID:      cmd.GoalID,
		Title:       h.goalService.SanitizeGoalTitle(cmd.Title),
		Description: h.goalService.SanitizeGoalDescription(cmd.Description),
		TargetDate:  cmd.TargetDate,
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   now,
	}
	
	// Validate milestone
	if err := h.goalService.ValidateMilestoneCreation(milestone); err != nil {
		return nil, fmt.Errorf("milestone validation failed: %w", err)
	}
	
	// Save milestone
	if err := h.milestoneRepo.Create(ctx, milestone); err != nil {
		return nil, fmt.Errorf("failed to create milestone: %w", err)
	}
	
	return &commands.CreateMilestoneResult{
		MilestoneID: milestone.ID,
		CreatedAt:   milestone.CreatedAt,
	}, nil
}

func (h *GoalHandler) HandleCompleteMilestone(ctx context.Context, cmd commands.CompleteMilestoneCommand) error {
	// Get milestone to check ownership
	milestone, err := h.milestoneRepo.GetByID(ctx, cmd.MilestoneID)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %w", err)
	}
	
	if milestone == nil {
		return fmt.Errorf("milestone not found")
	}
	
	// Get goal to check ownership
	goal, err := h.goalRepo.GetByID(ctx, milestone.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal.UserID != cmd.UserID {
		return fmt.Errorf("access denied: milestone belongs to different user")
	}
	
	// Mark milestone as completed
	if err := h.milestoneRepo.MarkCompleted(ctx, cmd.MilestoneID); err != nil {
		return fmt.Errorf("failed to complete milestone: %w", err)
	}
	
	// Recalculate goal progress
	tasks, err := h.taskRepo.GetByGoalID(ctx, milestone.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal tasks: %w", err)
	}
	
	milestones, err := h.milestoneRepo.GetByGoalID(ctx, milestone.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal milestones: %w", err)
	}
	
	newProgress := h.goalService.CalculateGoalProgress(tasks, milestones)
	if err := h.goalRepo.UpdateProgress(ctx, milestone.GoalID, newProgress); err != nil {
		return fmt.Errorf("failed to update goal progress: %w", err)
	}
	
	return nil
}

// Query Handlers

func (h *GoalHandler) HandleGetGoalByID(ctx context.Context, query queries.GetGoalByIDQuery) (*queries.GetGoalResult, error) {
	goal, err := h.goalRepo.GetByID(ctx, query.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}
	
	// Check ownership
	if goal.UserID != query.UserID {
		return nil, fmt.Errorf("access denied: goal belongs to different user")
	}
	
	return &queries.GetGoalResult{
		Goal: goal,
	}, nil
}

func (h *GoalHandler) HandleGetGoalsByUserID(ctx context.Context, query queries.GetGoalsByUserIDQuery) (*queries.GetGoalsResult, error) {
	goals, totalCount, err := h.goalRepo.GetByUserIDPaginated(ctx, query.UserID, query.Offset, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user goals: %w", err)
	}
	
	return &queries.GetGoalsResult{
		Goals:      goals,
		TotalCount: totalCount,
		Offset:     query.Offset,
		Limit:      query.Limit,
	}, nil
}

func (h *GoalHandler) HandleGetGoalsByStatus(ctx context.Context, query queries.GetGoalsByStatusQuery) (*queries.GetGoalsResult, error) {
	goals, err := h.goalRepo.GetByUserIDAndStatus(ctx, query.UserID, query.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals by status: %w", err)
	}
	
	return &queries.GetGoalsResult{
		Goals:      goals,
		TotalCount: int64(len(goals)),
		Offset:     query.Offset,
		Limit:      query.Limit,
	}, nil
}

func (h *GoalHandler) HandleGetTasksByGoalID(ctx context.Context, query queries.GetTasksByGoalIDQuery) (*queries.GetTasksResult, error) {
	// Check goal ownership
	goal, err := h.goalRepo.GetByID(ctx, query.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}
	
	if goal.UserID != query.UserID {
		return nil, fmt.Errorf("access denied: goal belongs to different user")
	}
	
	tasks, err := h.taskRepo.GetByGoalID(ctx, query.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal tasks: %w", err)
	}
	
	return &queries.GetTasksResult{
		Tasks:      tasks,
		TotalCount: int64(len(tasks)),
		Offset:     0,
		Limit:      len(tasks),
	}, nil
}

func (h *GoalHandler) HandleGetMilestonesByGoalID(ctx context.Context, query queries.GetMilestonesByGoalIDQuery) (*queries.GetMilestonesResult, error) {
	// Check goal ownership
	goal, err := h.goalRepo.GetByID(ctx, query.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}
	
	if goal.UserID != query.UserID {
		return nil, fmt.Errorf("access denied: goal belongs to different user")
	}
	
	milestones, err := h.milestoneRepo.GetByGoalID(ctx, query.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal milestones: %w", err)
	}
	
	return &queries.GetMilestonesResult{
		Milestones: milestones,
	}, nil
}