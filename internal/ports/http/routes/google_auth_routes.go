package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

func SetupGoogleAuthRoutes(
	router *gin.RouterGroup,
	googleAuthHandler *handlers.GoogleAuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	googleGroup := router.Group("/google")
	googleGroup.Use(authMiddleware.RequireAuth())

	// OAuth2 flow endpoints
	googleGroup.GET("/auth-url", googleAuthHandler.GetAuthURL)
	googleGroup.POST("/callback", googleAuthHandler.HandleCallback)
	
	// Integration management
	googleGroup.GET("/integration", googleAuthHandler.GetIntegration)
	googleGroup.DELETE("/integration", googleAuthHandler.DisconnectIntegration)
	
	// Calendar management
	googleGroup.GET("/calendars", googleAuthHandler.GetCalendars)
}