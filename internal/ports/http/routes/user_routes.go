package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

// SetupUserRoutes sets up all user-related routes
func SetupUserRoutes(
	router *gin.RouterGroup,
	userHandler *handlers.UserHTTPHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Public routes (no authentication required)
	auth := router.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.RefreshToken)
	}
	
	// Protected routes (authentication required)
	users := router.Group("/users")
	users.Use(authMiddleware.RequireAuth())
	{
		// Current user profile
		users.GET("/me", userHandler.GetProfile)
		users.PUT("/me", userHandler.UpdateProfile)
		users.DELETE("/me", userHandler.DeleteAccount)
	}
}