package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	appHandlers "github.com/andranikuz/smart-goal-calendar/internal/application/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/application/commands"
	"github.com/andranikuz/smart-goal-calendar/internal/application/queries"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

type EventHTTPHandler struct {
	eventHandler *appHandlers.EventHandler
}

func NewEventHTTPHandler(eventHandler *appHandlers.EventHandler) *EventHTTPHandler {
	return &EventHTTPHandler{
		eventHandler: eventHandler,
	}
}

// Request/Response models

type CreateEventRequest struct {
	GoalID         string                  `json:"goal_id,omitempty"`
	Title          string                  `json:"title" binding:"required,min=2,max=255"`
	Description    string                  `json:"description" binding:"max=1000"`
	StartTime      time.Time               `json:"start_time" binding:"required"`
	EndTime        time.Time               `json:"end_time" binding:"required"`
	Timezone       string                  `json:"timezone"`
	Recurrence     *entities.RecurrenceRule `json:"recurrence,omitempty"`
	Location       string                  `json:"location" binding:"max=255"`
	Attendees      []entities.Attendee     `json:"attendees"`
	Status         entities.EventStatus    `json:"status"`
	ExternalID     string                  `json:"external_id" binding:"max=255"`
	ExternalSource string                  `json:"external_source" binding:"max=50"`
}

type UpdateEventRequest struct {
	GoalID         *string                  `json:"goal_id,omitempty"`
	Title          *string                  `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	Description    *string                  `json:"description,omitempty" binding:"omitempty,max=1000"`
	StartTime      *time.Time               `json:"start_time,omitempty"`
	EndTime        *time.Time               `json:"end_time,omitempty"`
	Timezone       *string                  `json:"timezone,omitempty"`
	Recurrence     *entities.RecurrenceRule `json:"recurrence,omitempty"`
	Location       *string                  `json:"location,omitempty" binding:"omitempty,max=255"`
	Attendees      *[]entities.Attendee     `json:"attendees,omitempty"`
	Status         *entities.EventStatus    `json:"status,omitempty"`
	ExternalID     *string                  `json:"external_id,omitempty" binding:"omitempty,max=255"`
	ExternalSource *string                  `json:"external_source,omitempty" binding:"omitempty,max=50"`
}

type MoveEventRequest struct {
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

type DuplicateEventRequest struct {
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

type ChangeEventStatusRequest struct {
	Status entities.EventStatus `json:"status" binding:"required"`
}

type LinkEventToGoalRequest struct {
	GoalID string `json:"goal_id" binding:"required"`
}

type EventResponse struct {
	ID             entities.EventID         `json:"id"`
	UserID         entities.UserID          `json:"user_id"`
	GoalID         *entities.GoalID         `json:"goal_id"`
	Title          string                   `json:"title"`
	Description    string                   `json:"description"`
	StartTime      time.Time                `json:"start_time"`
	EndTime        time.Time                `json:"end_time"`
	Timezone       string                   `json:"timezone"`
	Recurrence     *entities.RecurrenceRule `json:"recurrence"`
	Location       string                   `json:"location"`
	Attendees      []entities.Attendee      `json:"attendees"`
	Status         entities.EventStatus     `json:"status"`
	ExternalID     string                   `json:"external_id"`
	ExternalSource string                   `json:"external_source"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

// Event CRUD operations

func (h *EventHTTPHandler) CreateEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Convert goal ID if provided
	var goalID *entities.GoalID
	if req.GoalID != "" {
		gID := entities.GoalID(req.GoalID)
		goalID = &gID
	}
	
	// Set default timezone if not provided
	if req.Timezone == "" {
		req.Timezone = "UTC"
	}
	
	// Set default status if not provided
	if req.Status == "" {
		req.Status = entities.EventStatusConfirmed
	}
	
	// Create command
	cmd := commands.CreateEventCommand{
		UserID:         userID,
		GoalID:         goalID,
		Title:          req.Title,
		Description:    req.Description,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Timezone:       req.Timezone,
		Recurrence:     req.Recurrence,
		Location:       req.Location,
		Attendees:      req.Attendees,
		Status:         req.Status,
		ExternalID:     req.ExternalID,
		ExternalSource: req.ExternalSource,
	}
	
	// Execute command
	result, err := h.eventHandler.HandleCreateEvent(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "event_creation_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Event created successfully",
		"event_id":   result.EventID,
		"created_at": result.CreatedAt,
	})
}

func (h *EventHTTPHandler) GetEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	eventID := entities.EventID(c.Param("id"))
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Event ID is required",
		})
		return
	}
	
	// Create query
	query := queries.GetEventByIDQuery{
		EventID: eventID,
		UserID:  userID,
	}
	
	// Execute query
	result, err := h.eventHandler.HandleGetEventByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "event_not_found",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"event": h.mapEventToResponse(result.Event),
	})
}

func (h *EventHTTPHandler) GetEvents(c *gin.Context) {
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
	query := queries.GetEventsByUserIDQuery{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	}
	
	// Execute query
	result, err := h.eventHandler.HandleGetEventsByUserID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "events_retrieval_failed",
			"message": err.Error(),
		})
		return
	}
	
	// Map events to response format
	events := make([]EventResponse, len(result.Events))
	for i, event := range result.Events {
		events[i] = h.mapEventToResponse(event)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"events":      events,
		"total_count": result.TotalCount,
		"offset":      result.Offset,
		"limit":       result.Limit,
	})
}

func (h *EventHTTPHandler) GetEventsByTimeRange(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	// Parse time range parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	
	if startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameters",
			"message": "start_time and end_time are required",
		})
		return
	}
	
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_start_time",
			"message": "start_time must be in RFC3339 format",
		})
		return
	}
	
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_end_time",
			"message": "end_time must be in RFC3339 format",
		})
		return
	}
	
	// Create query
	query := queries.GetEventsByTimeRangeQuery{
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime,
	}
	
	// Execute query
	result, err := h.eventHandler.HandleGetEventsByTimeRange(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "events_retrieval_failed",
			"message": err.Error(),
		})
		return
	}
	
	// Map events to response format
	events := make([]EventResponse, len(result.Events))
	for i, event := range result.Events {
		events[i] = h.mapEventToResponse(event)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}

func (h *EventHTTPHandler) UpdateEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	eventID := entities.EventID(c.Param("id"))
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Event ID is required",
		})
		return
	}
	
	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Convert goal ID if provided
	var goalID *entities.GoalID
	if req.GoalID != nil && *req.GoalID != "" {
		gID := entities.GoalID(*req.GoalID)
		goalID = &gID
	}
	
	// Create command
	cmd := commands.UpdateEventCommand{
		EventID:        eventID,
		UserID:         userID,
		GoalID:         goalID,
		Title:          req.Title,
		Description:    req.Description,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Timezone:       req.Timezone,
		Recurrence:     req.Recurrence,
		Location:       req.Location,
		Attendees:      req.Attendees,
		Status:         req.Status,
		ExternalID:     req.ExternalID,
		ExternalSource: req.ExternalSource,
	}
	
	// Execute command
	result, err := h.eventHandler.HandleUpdateEvent(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "event_update_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Event updated successfully",
		"updated_at": result.UpdatedAt,
	})
}

func (h *EventHTTPHandler) DeleteEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	eventID := entities.EventID(c.Param("id"))
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Event ID is required",
		})
		return
	}
	
	// Create command
	cmd := commands.DeleteEventCommand{
		EventID: eventID,
		UserID:  userID,
	}
	
	// Execute command
	if err := h.eventHandler.HandleDeleteEvent(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "event_deletion_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
}

// Event actions

func (h *EventHTTPHandler) MoveEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	eventID := entities.EventID(c.Param("id"))
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Event ID is required",
		})
		return
	}
	
	var req MoveEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Create command
	cmd := commands.MoveEventCommand{
		EventID:   eventID,
		UserID:    userID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	
	// Execute command
	result, err := h.eventHandler.HandleMoveEvent(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "event_move_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Event moved successfully",
		"updated_at": result.UpdatedAt,
	})
}

func (h *EventHTTPHandler) DuplicateEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	eventID := entities.EventID(c.Param("id"))
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Event ID is required",
		})
		return
	}
	
	var req DuplicateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Create command
	cmd := commands.DuplicateEventCommand{
		EventID:   eventID,
		UserID:    userID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	
	// Execute command
	result, err := h.eventHandler.HandleDuplicateEvent(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "event_duplication_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Event duplicated successfully",
		"event_id":   result.EventID,
		"created_at": result.CreatedAt,
	})
}

func (h *EventHTTPHandler) ChangeEventStatus(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	eventID := entities.EventID(c.Param("id"))
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Event ID is required",
		})
		return
	}
	
	var req ChangeEventStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Create command
	cmd := commands.ChangeEventStatusCommand{
		EventID: eventID,
		UserID:  userID,
		Status:  req.Status,
	}
	
	// Execute command
	result, err := h.eventHandler.HandleChangeEventStatus(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "status_change_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Event status changed successfully",
		"updated_at": result.UpdatedAt,
	})
}

func (h *EventHTTPHandler) LinkEventToGoal(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	eventID := entities.EventID(c.Param("id"))
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Event ID is required",
		})
		return
	}
	
	var req LinkEventToGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Create command
	cmd := commands.LinkEventToGoalCommand{
		EventID: eventID,
		UserID:  userID,
		GoalID:  entities.GoalID(req.GoalID),
	}
	
	// Execute command
	result, err := h.eventHandler.HandleLinkEventToGoal(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "link_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Event linked to goal successfully",
		"updated_at": result.UpdatedAt,
	})
}

func (h *EventHTTPHandler) UnlinkEventFromGoal(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	eventID := entities.EventID(c.Param("id"))
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Event ID is required",
		})
		return
	}
	
	// Create command
	cmd := commands.UnlinkEventFromGoalCommand{
		EventID: eventID,
		UserID:  userID,
	}
	
	// Execute command
	result, err := h.eventHandler.HandleUnlinkEventFromGoal(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "unlink_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Event unlinked from goal successfully",
		"updated_at": result.UpdatedAt,
	})
}

// Additional event endpoints

func (h *EventHTTPHandler) GetUpcomingEvents(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	// Parse limit parameter
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 50 {
		limit = 50
	}
	
	// Create query
	query := queries.GetUpcomingEventsQuery{
		UserID: userID,
		Limit:  limit,
	}
	
	// Execute query
	result, err := h.eventHandler.HandleGetUpcomingEvents(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "events_retrieval_failed",
			"message": err.Error(),
		})
		return
	}
	
	// Map events to response format
	events := make([]EventResponse, len(result.Events))
	for i, event := range result.Events {
		events[i] = h.mapEventToResponse(event)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}

func (h *EventHTTPHandler) GetTodayEvents(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	// Get timezone from query parameter or default to UTC
	timezone := c.DefaultQuery("timezone", "UTC")
	
	// Create query
	query := queries.GetTodayEventsQuery{
		UserID:   userID,
		Timezone: timezone,
	}
	
	// Execute query
	result, err := h.eventHandler.HandleGetTodayEvents(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "events_retrieval_failed",
			"message": err.Error(),
		})
		return
	}
	
	// Map events to response format
	events := make([]EventResponse, len(result.Events))
	for i, event := range result.Events {
		events[i] = h.mapEventToResponse(event)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}

func (h *EventHTTPHandler) SearchEvents(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	// Get search query
	searchQuery := c.Query("q")
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "Search query 'q' is required",
		})
		return
	}
	
	// Parse limit parameter
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 50 {
		limit = 50
	}
	
	// Create query
	query := queries.SearchEventsQuery{
		UserID: userID,
		Query:  searchQuery,
		Limit:  limit,
	}
	
	// Execute query
	result, err := h.eventHandler.HandleSearchEvents(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "search_failed",
			"message": err.Error(),
		})
		return
	}
	
	// Map events to response format
	events := make([]EventResponse, len(result.Events))
	for i, event := range result.Events {
		events[i] = h.mapEventToResponse(event)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"query":  searchQuery,
	})
}

func (h *EventHTTPHandler) CheckEventConflict(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	// Parse time range parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	excludeEventIDStr := c.Query("exclude_event_id")
	
	if startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameters",
			"message": "start_time and end_time are required",
		})
		return
	}
	
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_start_time",
			"message": "start_time must be in RFC3339 format",
		})
		return
	}
	
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_end_time",
			"message": "end_time must be in RFC3339 format",
		})
		return
	}
	
	var excludeEventID *entities.EventID
	if excludeEventIDStr != "" {
		eID := entities.EventID(excludeEventIDStr)
		excludeEventID = &eID
	}
	
	// Create query
	query := queries.CheckEventConflictQuery{
		UserID:          userID,
		StartTime:       startTime,
		EndTime:         endTime,
		ExcludeEventID:  excludeEventID,
	}
	
	// Execute query
	result, err := h.eventHandler.HandleCheckEventConflict(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "conflict_check_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"has_conflict": result.HasConflict,
	})
}

// Helper methods

func (h *EventHTTPHandler) mapEventToResponse(event *entities.Event) EventResponse {
	return EventResponse{
		ID:             event.ID,
		UserID:         event.UserID,
		GoalID:         event.GoalID,
		Title:          event.Title,
		Description:    event.Description,
		StartTime:      event.StartTime,
		EndTime:        event.EndTime,
		Timezone:       event.Timezone,
		Recurrence:     event.Recurrence,
		Location:       event.Location,
		Attendees:      event.Attendees,
		Status:         event.Status,
		ExternalID:     event.ExternalID,
		ExternalSource: event.ExternalSource,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}