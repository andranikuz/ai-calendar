package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

func SetupEventRoutes(
	router *gin.RouterGroup,
	eventHandler *handlers.EventHTTPHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Event routes group
	events := router.Group("/events")
	events.Use(authMiddleware.RequireAuth())

	// Event CRUD operations
	events.POST("", eventHandler.CreateEvent)                    // Create event
	events.GET("", eventHandler.GetEvents)                       // Get user's events (paginated)
	events.GET("/search", eventHandler.SearchEvents)             // Search events
	events.GET("/upcoming", eventHandler.GetUpcomingEvents)      // Get upcoming events
	events.GET("/today", eventHandler.GetTodayEvents)           // Get today's events
	events.GET("/time-range", eventHandler.GetEventsByTimeRange) // Get events by time range
	events.GET("/conflict-check", eventHandler.CheckEventConflict) // Check for conflicts
	events.GET("/:id", eventHandler.GetEvent)                   // Get specific event
	events.PUT("/:id", eventHandler.UpdateEvent)                // Update event
	events.DELETE("/:id", eventHandler.DeleteEvent)             // Delete event

	// Event actions
	events.POST("/:id/move", eventHandler.MoveEvent)            // Move event to new time
	events.POST("/:id/duplicate", eventHandler.DuplicateEvent) // Duplicate event
	events.POST("/:id/status", eventHandler.ChangeEventStatus) // Change event status
	events.POST("/:id/link-goal", eventHandler.LinkEventToGoal) // Link event to goal
	events.POST("/:id/unlink-goal", eventHandler.UnlinkEventFromGoal) // Unlink event from goal
}