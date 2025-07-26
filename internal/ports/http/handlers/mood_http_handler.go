package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/application/commands"
	"github.com/andranikuz/smart-goal-calendar/internal/application/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/application/queries"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

type MoodHTTPHandler struct {
	moodHandler *handlers.MoodHandler
}

func NewMoodHTTPHandler(moodHandler *handlers.MoodHandler) *MoodHTTPHandler {
	return &MoodHTTPHandler{
		moodHandler: moodHandler,
	}
}

type CreateMoodRequest struct {
	Date  string                `json:"date" binding:"required"`
	Level entities.MoodLevel    `json:"level" binding:"required,min=1,max=5"`
	Notes string                `json:"notes"`
	Tags  []entities.MoodTag    `json:"tags"`
}

type UpdateMoodRequest struct {
	Level entities.MoodLevel    `json:"level" binding:"required,min=1,max=5"`
	Notes string                `json:"notes"`
	Tags  []entities.MoodTag    `json:"tags"`
}

type UpsertMoodRequest struct {
	Date  string                `json:"date" binding:"required"`
	Level entities.MoodLevel    `json:"level" binding:"required,min=1,max=5"`
	Notes string                `json:"notes"`
	Tags  []entities.MoodTag    `json:"tags"`
}

type MoodResponse struct {
	ID         string                `json:"id"`
	UserID     string                `json:"user_id"`
	Date       string                `json:"date"`
	Level      entities.MoodLevel    `json:"level"`
	LevelText  string                `json:"level_text"`
	LevelEmoji string                `json:"level_emoji"`
	Notes      string                `json:"notes"`
	Tags       []entities.MoodTag    `json:"tags"`
	RecordedAt time.Time             `json:"recorded_at"`
}

func (h *MoodHTTPHandler) CreateMood(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateMoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	cmd := commands.CreateMoodCommand{
		UserID: userID,
		Date:   date,
		Level:  req.Level,
		Notes:  req.Notes,
		Tags:   req.Tags,
	}

	mood, err := h.moodHandler.CreateMood(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.toMoodResponse(mood)
	c.JSON(http.StatusCreated, response)
}

func (h *MoodHTTPHandler) GetMoods(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	query := queries.GetMoodsByUserIDQuery{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}

	moods, err := h.moodHandler.GetMoodsByUserID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]MoodResponse, len(moods))
	for i, mood := range moods {
		responses[i] = h.toMoodResponse(mood)
	}

	c.JSON(http.StatusOK, gin.H{
		"moods":  responses,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}

func (h *MoodHTTPHandler) GetMood(c *gin.Context) {
	moodID := c.Param("id")
	if moodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mood ID is required"})
		return
	}

	query := queries.GetMoodByIDQuery{
		ID: entities.MoodID(moodID),
	}

	mood, err := h.moodHandler.GetMoodByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if mood == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mood not found"})
		return
	}

	response := h.toMoodResponse(mood)
	c.JSON(http.StatusOK, response)
}

func (h *MoodHTTPHandler) UpdateMood(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	moodID := c.Param("id")
	if moodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mood ID is required"})
		return
	}

	var req UpdateMoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.UpdateMoodCommand{
		ID:     entities.MoodID(moodID),
		UserID: userID,
		Level:  req.Level,
		Notes:  req.Notes,
		Tags:   req.Tags,
	}

	mood, err := h.moodHandler.UpdateMood(c.Request.Context(), cmd)
	if err != nil {
		if err.Error() == "mood not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mood not found"})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this mood"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.toMoodResponse(mood)
	c.JSON(http.StatusOK, response)
}

func (h *MoodHTTPHandler) DeleteMood(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	moodID := c.Param("id")
	if moodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mood ID is required"})
		return
	}

	cmd := commands.DeleteMoodCommand{
		ID:     entities.MoodID(moodID),
		UserID: userID,
	}

	err := h.moodHandler.DeleteMood(c.Request.Context(), cmd)
	if err != nil {
		if err.Error() == "mood not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mood not found"})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this mood"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mood deleted successfully"})
}

func (h *MoodHTTPHandler) GetMoodByDate(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	query := queries.GetMoodByUserIDAndDateQuery{
		UserID: userID,
		Date:   date,
	}

	mood, err := h.moodHandler.GetMoodByUserIDAndDate(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if mood == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No mood found for this date"})
		return
	}

	response := h.toMoodResponse(mood)
	c.JSON(http.StatusOK, response)
}

func (h *MoodHTTPHandler) GetMoodsByDateRange(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both start and end date parameters are required"})
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	query := queries.GetMoodsByDateRangeQuery{
		UserID: userID,
		Start:  start,
		End:    end,
	}

	moods, err := h.moodHandler.GetMoodsByDateRange(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]MoodResponse, len(moods))
	for i, mood := range moods {
		responses[i] = h.toMoodResponse(mood)
	}

	c.JSON(http.StatusOK, gin.H{
		"moods": responses,
		"start": startStr,
		"end":   endStr,
		"count": len(responses),
	})
}

func (h *MoodHTTPHandler) GetLatestMood(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	query := queries.GetLatestMoodQuery{
		UserID: userID,
	}

	mood, err := h.moodHandler.GetLatestMood(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if mood == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No mood entries found"})
		return
	}

	response := h.toMoodResponse(mood)
	c.JSON(http.StatusOK, response)
}

func (h *MoodHTTPHandler) UpsertMoodByDate(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req UpsertMoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	cmd := commands.UpsertMoodByDateCommand{
		UserID: userID,
		Date:   date,
		Level:  req.Level,
		Notes:  req.Notes,
		Tags:   req.Tags,
	}

	mood, err := h.moodHandler.UpsertMoodByDate(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.toMoodResponse(mood)
	c.JSON(http.StatusOK, response)
}

func (h *MoodHTTPHandler) GetMoodStats(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	startStr := c.DefaultQuery("start", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endStr := c.DefaultQuery("end", time.Now().Format("2006-01-02"))

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	query := queries.GetMoodStatsQuery{
		UserID: userID,
		Start:  start,
		End:    end,
	}

	stats, err := h.moodHandler.GetMoodStats(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *MoodHTTPHandler) GetMoodTrends(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}
	if days > 365 {
		days = 365
	}

	query := queries.GetMoodTrendsQuery{
		UserID: userID,
		Days:   days,
	}

	trends, err := h.moodHandler.GetMoodTrends(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trends": trends,
		"days":   days,
		"count":  len(trends),
	})
}

func (h *MoodHTTPHandler) toMoodResponse(mood *entities.Mood) MoodResponse {
	return MoodResponse{
		ID:         string(mood.ID),
		UserID:     string(mood.UserID),
		Date:       mood.Date.Format("2006-01-02"),
		Level:      mood.Level,
		LevelText:  mood.Level.String(),
		LevelEmoji: mood.Level.Emoji(),
		Notes:      mood.Notes,
		Tags:       mood.Tags,
		RecordedAt: mood.RecordedAt,
	}
}