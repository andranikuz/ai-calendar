package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

func SetupGoalRoutes(
	router *gin.RouterGroup,
	goalHandler *handlers.GoalHTTPHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Goal routes group
	goals := router.Group("/goals")
	goals.Use(authMiddleware.RequireAuth())

	// Goal CRUD operations
	goals.POST("", goalHandler.CreateGoal)                  // Create goal
	goals.GET("", goalHandler.GetGoals)                     // Get user's goals (paginated)
	goals.GET("/:id", goalHandler.GetGoal)                  // Get specific goal
	goals.PUT("/:id", goalHandler.UpdateGoal)               // Update goal
	goals.DELETE("/:id", goalHandler.DeleteGoal)            // Delete goal

	// Task management within goals
	goals.POST("/:id/tasks", goalHandler.CreateTask)        // Create task for goal
	goals.GET("/:id/tasks", goalHandler.GetTasks)           // Get tasks for goal
	goals.POST("/tasks/:taskId/complete", goalHandler.CompleteTask) // Complete task

	// Milestone management within goals
	goals.POST("/:id/milestones", goalHandler.CreateMilestone)     // Create milestone for goal
	goals.GET("/:id/milestones", goalHandler.GetMilestones)        // Get milestones for goal
	goals.POST("/milestones/:milestoneId/complete", goalHandler.CompleteMilestone) // Complete milestone
}