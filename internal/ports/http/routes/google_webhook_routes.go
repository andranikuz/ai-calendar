package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

func RegisterGoogleWebhookRoutes(r *gin.RouterGroup, webhookHandler *handlers.GoogleWebhookHandler, authMiddleware *middleware.AuthMiddleware) {
	webhookGroup := r.Group("/google")
	{
		// Webhook endpoint - no auth required as Google sends requests here
		webhookGroup.POST("/webhook", webhookHandler.HandleCalendarWebhook)
		
		// Setup webhook endpoint - requires authentication
		webhookGroup.POST("/webhook/setup", authMiddleware.RequireAuth(), webhookHandler.SetupWebhook)
	}
}