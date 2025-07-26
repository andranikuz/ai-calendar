package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	appHandlers "github.com/andranikuz/smart-goal-calendar/internal/application/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/application/commands"
	"github.com/andranikuz/smart-goal-calendar/internal/application/queries"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
	"time"
)

type GoalHTTPHandler struct {
	goalHandler *appHandlers.GoalHandler
}

func NewGoalHTTPHandler(goalHandler *appHandlers.GoalHandler) *GoalHTTPHandler {
	return &GoalHTTPHandler{
		goalHandler: goalHandler,
	}
}

// Request/Response models

type CreateGoalRequest struct {
	Title       string                `json:"title" binding:"required,min=2,max=255"`
	Description string                `json:"description" binding:"max=1000"`
	Category    entities.GoalCategory `json:"category" binding:"required"`
	Priority    entities.Priority     `json:"priority" binding:"required"`
	Deadline    *time.Time            `json:"deadline,omitempty"`
}

type UpdateGoalRequest struct {
	Title       *string                `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string                `json:"description,omitempty" binding:"omitempty,max=1000"`
	Category    *entities.GoalCategory `json:"category,omitempty"`
	Priority    *entities.Priority     `json:"priority,omitempty"`
	Status      *entities.GoalStatus   `json:"status,omitempty"`
	Progress    *int                   `json:"progress,omitempty" binding:"omitempty,min=0,max=100"`
	Deadline    *time.Time             `json:"deadline,omitempty"`
}

type CreateTaskRequest struct {
	Title             string         `json:"title" binding:"required,min=2,max=255"`
	Description       string         `json:"description" binding:"max=1000"`
	Priority          entities.Priority `json:"priority" binding:"required"`
	EstimatedDuration int            `json:"estimated_duration" binding:"min=1"` // minutes
	DueDate           *time.Time     `json:"due_date,omitempty"`
}

type CreateMilestoneRequest struct {
	Title       string    `json:"title" binding:"required,min=2,max=255"`
	Description string    `json:"description" binding:"max=1000"`
	TargetDate  time.Time `json:"target_date" binding:"required"`
}

type GoalResponse struct {
	ID          entities.GoalID       `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Category    entities.GoalCategory `json:"category"`
	Priority    entities.Priority     `json:"priority"`
	Status      entities.GoalStatus   `json:"status"`
	Progress    int                   `json:"progress"`
	Deadline    *time.Time            `json:"deadline"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// Goal CRUD operations

func (h *GoalHTTPHandler) CreateGoal(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Create command
	cmd := commands.CreateGoalCommand{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Priority:    req.Priority,
		Deadline:    req.Deadline,
	}
	
	// Execute command
	result, err := h.goalHandler.HandleCreateGoal(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "goal_creation_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Goal created successfully",
		"goal_id":    result.GoalID,
		"created_at": result.CreatedAt,
	})
}

func (h *GoalHTTPHandler) GetGoal(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	goalID := entities.GoalID(c.Param("id"))
	if goalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Goal ID is required",
		})
		return
	}
	
	// Create query
	query := queries.GetGoalByIDQuery{
		GoalID: goalID,
		UserID: userID,
	}
	
	// Execute query
	result, err := h.goalHandler.HandleGetGoalByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "goal_not_found",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"goal": h.mapGoalToResponse(result.Goal),
	})
}

func (h *GoalHTTPHandler) GetGoals(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	// Parse pagination parameters
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	// Limit the maximum number of results
	if limit > 100 {
		limit = 100
	}
	
	// Create query
	query := queries.GetGoalsByUserIDQuery{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	}
	
	// Execute query
	result, err := h.goalHandler.HandleGetGoalsByUserID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "goals_retrieval_failed",
			"message": err.Error(),
		})
		return
	}
	
	// Map goals to response format
	goals := make([]GoalResponse, len(result.Goals))
	for i, goal := range result.Goals {
		goals[i] = h.mapGoalToResponse(goal)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"goals":       goals,
		"total_count": result.TotalCount,
		"offset":      result.Offset,
		"limit":       result.Limit,
	})
}

func (h *GoalHTTPHandler) UpdateGoal(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	goalID := entities.GoalID(c.Param("id"))
	if goalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Goal ID is required",
		})
		return
	}
	
	var req UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Create command
	cmd := commands.UpdateGoalCommand{
		GoalID:      goalID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Priority:    req.Priority,
		Status:      req.Status,
		Progress:    req.Progress,
		Deadline:    req.Deadline,
	}
	
	// Execute command
	result, err := h.goalHandler.HandleUpdateGoal(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "goal_update_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Goal updated successfully",
		"updated_at": result.UpdatedAt,
	})
}

func (h *GoalHTTPHandler) DeleteGoal(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	goalID := entities.GoalID(c.Param("id"))
	if goalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Goal ID is required",
		})
		return
	}
	
	// Create command
	cmd := commands.DeleteGoalCommand{
		GoalID: goalID,
		UserID: userID,
	}
	
	// Execute command
	if err := h.goalHandler.HandleDeleteGoal(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "goal_deletion_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Goal deleted successfully",
	})
}

// Task operations

func (h *GoalHTTPHandler) CreateTask(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	goalID := entities.GoalID(c.Param("id"))
	if goalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Goal ID is required",
		})
		return
	}
	
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Create command
	cmd := commands.CreateTaskCommand{
		GoalID:            goalID,
		UserID:            userID,
		Title:             req.Title,
		Description:       req.Description,
		Priority:          req.Priority,
		EstimatedDuration: req.EstimatedDuration,
		DueDate:           req.DueDate,
	}
	
	// Execute command
	result, err := h.goalHandler.HandleCreateTask(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "task_creation_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Task created successfully",
		"task_id":    result.TaskID,
		"created_at": result.CreatedAt,
	})
}

func (h *GoalHTTPHandler) GetTasks(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	goalID := entities.GoalID(c.Param("id"))
	if goalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Goal ID is required",
		})
		return
	}
	
	// Create query
	query := queries.GetTasksByGoalIDQuery{
		GoalID: goalID,
		UserID: userID,
	}
	
	// Execute query
	result, err := h.goalHandler.HandleGetTasksByGoalID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "tasks_retrieval_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"tasks":       result.Tasks,
		"total_count": result.TotalCount,
	})
}

func (h *GoalHTTPHandler) CompleteTask(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	taskID := entities.TaskID(c.Param("taskId"))
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Task ID is required",
		})
		return
	}
	
	// Create command
	cmd := commands.CompleteTaskCommand{
		TaskID: taskID,
		UserID: userID,
	}
	
	// Execute command
	if err := h.goalHandler.HandleCompleteTask(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "task_completion_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Task completed successfully",
	})
}

// Milestone operations

func (h *GoalHTTPHandler) CreateMilestone(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	goalID := entities.GoalID(c.Param("id"))
	if goalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Goal ID is required",
		})
		return
	}
	
	var req CreateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Create command
	cmd := commands.CreateMilestoneCommand{
		GoalID:      goalID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		TargetDate:  req.TargetDate,
	}
	
	// Execute command
	result, err := h.goalHandler.HandleCreateMilestone(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "milestone_creation_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message":      "Milestone created successfully",
		"milestone_id": result.MilestoneID,
		"created_at":   result.CreatedAt,
	})
}

func (h *GoalHTTPHandler) GetMilestones(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	goalID := entities.GoalID(c.Param("id"))
	if goalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Goal ID is required",
		})
		return
	}
	
	// Create query
	query := queries.GetMilestonesByGoalIDQuery{
		GoalID: goalID,
		UserID: userID,
	}
	
	// Execute query
	result, err := h.goalHandler.HandleGetMilestonesByGoalID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "milestones_retrieval_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"milestones": result.Milestones,
	})
}

func (h *GoalHTTPHandler) CompleteMilestone(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	milestoneID := entities.MilestoneID(c.Param("milestoneId"))
	if milestoneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Milestone ID is required",
		})
		return
	}
	
	// Create command
	cmd := commands.CompleteMilestoneCommand{
		MilestoneID: milestoneID,
		UserID:      userID,
	}
	
	// Execute command
	if err := h.goalHandler.HandleCompleteMilestone(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "milestone_completion_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Milestone completed successfully",
	})
}

// Helper methods

func (h *GoalHTTPHandler) mapGoalToResponse(goal *entities.Goal) GoalResponse {
	return GoalResponse{
		ID:          goal.ID,
		Title:       goal.Title,
		Description: goal.Description,
		Category:    goal.Category,
		Priority:    goal.Priority,
		Status:      goal.Status,
		Progress:    goal.Progress,
		Deadline:    goal.Deadline,
		CreatedAt:   goal.CreatedAt,
		UpdatedAt:   goal.UpdatedAt,
	}
}