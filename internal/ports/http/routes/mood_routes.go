package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

func SetupMoodRoutes(
	router *gin.RouterGroup,
	moodHandler *handlers.MoodHTTPHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	moodGroup := router.Group("/moods")
	moodGroup.Use(authMiddleware.RequireAuth())

	// Basic CRUD operations
	moodGroup.POST("", moodHandler.CreateMood)
	moodGroup.GET("", moodHandler.GetMoods)
	moodGroup.GET("/:id", moodHandler.GetMood)
	moodGroup.PUT("/:id", moodHandler.UpdateMood)
	moodGroup.DELETE("/:id", moodHandler.DeleteMood)

	// Special operations
	moodGroup.GET("/by-date", moodHandler.GetMoodByDate)
	moodGroup.GET("/date-range", moodHandler.GetMoodsByDateRange)
	moodGroup.GET("/latest", moodHandler.GetLatestMood)
	moodGroup.POST("/upsert-by-date", moodHandler.UpsertMoodByDate)

	// Analytics
	moodGroup.GET("/stats", moodHandler.GetMoodStats)
	moodGroup.GET("/trends", moodHandler.GetMoodTrends)
}