package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

func RegisterSyncConflictRoutes(router *gin.RouterGroup, handler *handlers.SyncConflictHandler, authMiddleware gin.HandlerFunc) {
	conflicts := router.Group("/sync-conflicts")
	conflicts.Use(authMiddleware)
	{
		// Get all pending conflicts for user
		conflicts.GET("", handler.GetPendingConflicts)
		
		// Get conflict statistics
		conflicts.GET("/stats", handler.GetConflictStats)
		
		// Get specific conflict details
		conflicts.GET("/:id", handler.GetConflictDetails)
		
		// Resolve a specific conflict
		conflicts.POST("/:id/resolve", handler.ResolveConflict)
		
		// Bulk resolve multiple conflicts
		conflicts.POST("/bulk-resolve", handler.BulkResolveConflicts)
	}
}