package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

func SetupGoogleCalendarSyncRoutes(
	router *gin.RouterGroup,
	googleCalendarSyncHandler *handlers.GoogleCalendarSyncHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	syncGroup := router.Group("/google/calendar-syncs")
	syncGroup.Use(authMiddleware.RequireAuth())

	// Calendar sync configuration endpoints
	syncGroup.POST("", googleCalendarSyncHandler.CreateCalendarSync)
	syncGroup.GET("", googleCalendarSyncHandler.GetCalendarSyncs)
	syncGroup.PUT("/:id", googleCalendarSyncHandler.UpdateCalendarSync)
	syncGroup.DELETE("/:id", googleCalendarSyncHandler.DeleteCalendarSync)
	
	// Sync action endpoints
	syncGroup.POST("/:id/sync", googleCalendarSyncHandler.SyncNow)
}