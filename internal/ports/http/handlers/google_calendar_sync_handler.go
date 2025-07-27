package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/andranikuz/smart-goal-calendar/internal/adapters/google"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

type GoogleCalendarSyncHandler struct {
	oauth2Service          *google.OAuth2Service
	calendarService        *google.CalendarService
	googleIntegrationRepo  repositories.GoogleIntegrationRepository
	googleCalendarSyncRepo repositories.GoogleCalendarSyncRepository
	eventRepo              repositories.EventRepository
}

func NewGoogleCalendarSyncHandler(
	oauth2Service *google.OAuth2Service,
	calendarService *google.CalendarService,
	googleIntegrationRepo repositories.GoogleIntegrationRepository,
	googleCalendarSyncRepo repositories.GoogleCalendarSyncRepository,
	eventRepo repositories.EventRepository,
) *GoogleCalendarSyncHandler {
	return &GoogleCalendarSyncHandler{
		oauth2Service:          oauth2Service,
		calendarService:        calendarService,
		googleIntegrationRepo:  googleIntegrationRepo,
		googleCalendarSyncRepo: googleCalendarSyncRepo,
		eventRepo:              eventRepo,
	}
}

type CalendarSyncConfigRequest struct {
	CalendarID    string                         `json:"calendar_id" binding:"required"`
	CalendarName  string                         `json:"calendar_name" binding:"required"`
	SyncDirection entities.CalendarSyncDirection `json:"sync_direction" binding:"required"`
	SyncStatus    entities.CalendarSyncStatus    `json:"sync_status"`
	Settings      entities.CalendarSyncSettings  `json:"settings"`
}

type CalendarSyncConfigResponse struct {
	ID            string                         `json:"id"`
	CalendarID    string                         `json:"calendar_id"`
	CalendarName  string                         `json:"calendar_name"`
	SyncDirection entities.CalendarSyncDirection `json:"sync_direction"`
	SyncStatus    entities.CalendarSyncStatus    `json:"sync_status"`
	LastSyncAt    *time.Time                     `json:"last_sync_at"`
	LastSyncError string                         `json:"last_sync_error"`
	Settings      entities.CalendarSyncSettings  `json:"settings"`
	CreatedAt     time.Time                      `json:"created_at"`
	UpdatedAt     time.Time                      `json:"updated_at"`
}

// CreateCalendarSync creates a new calendar sync configuration
func (h *GoogleCalendarSyncHandler) CreateCalendarSync(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has Google integration
	integration, err := h.googleIntegrationRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Google integration"})
		return
	}
	if integration == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Google integration not found. Please connect Google account first"})
		return
	}

	var req CalendarSyncConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if sync already exists for this calendar
	existingSync, err := h.googleCalendarSyncRepo.GetByCalendarID(c.Request.Context(), req.CalendarID)
	if err == nil && existingSync != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Sync configuration already exists for this calendar"})
		return
	}

	// Set default settings if not provided
	if req.Settings.SyncInterval == 0 {
		req.Settings = entities.DefaultCalendarSyncSettings()
	}

	now := time.Now()
	sync := &entities.GoogleCalendarSync{
		ID:                  uuid.New().String(),
		UserID:              userID,
		GoogleIntegrationID: integration.ID,
		CalendarID:          req.CalendarID,
		CalendarName:        req.CalendarName,
		SyncDirection:       req.SyncDirection,
		SyncStatus:          entities.SyncStatusActive,
		Settings:            req.Settings,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := h.googleCalendarSyncRepo.Create(c.Request.Context(), sync); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create calendar sync configuration"})
		return
	}

	response := h.toSyncConfigResponse(sync)
	c.JSON(http.StatusCreated, response)
}

// GetCalendarSyncs returns all calendar sync configurations for the user
func (h *GoogleCalendarSyncHandler) GetCalendarSyncs(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	syncs, err := h.googleCalendarSyncRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendar sync configurations"})
		return
	}

	var response []CalendarSyncConfigResponse
	for _, sync := range syncs {
		response = append(response, h.toSyncConfigResponse(sync))
	}

	c.JSON(http.StatusOK, gin.H{"syncs": response})
}

// UpdateCalendarSync updates a calendar sync configuration
func (h *GoogleCalendarSyncHandler) UpdateCalendarSync(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	syncID := c.Param("id")
	if syncID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sync ID is required"})
		return
	}

	// Get existing sync
	sync, err := h.googleCalendarSyncRepo.GetByID(c.Request.Context(), syncID)
	if err != nil || sync == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar sync configuration not found"})
		return
	}

	// Verify ownership
	if sync.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this sync configuration"})
		return
	}

	var req CalendarSyncConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	sync.CalendarName = req.CalendarName
	sync.SyncDirection = req.SyncDirection
	if req.SyncStatus != "" {
		sync.SyncStatus = req.SyncStatus
	}
	sync.Settings = req.Settings
	sync.UpdatedAt = time.Now()

	if err := h.googleCalendarSyncRepo.Update(c.Request.Context(), sync); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update calendar sync configuration"})
		return
	}

	response := h.toSyncConfigResponse(sync)
	c.JSON(http.StatusOK, response)
}

// DeleteCalendarSync deletes a calendar sync configuration
func (h *GoogleCalendarSyncHandler) DeleteCalendarSync(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	syncID := c.Param("id")
	if syncID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sync ID is required"})
		return
	}

	// Get existing sync
	sync, err := h.googleCalendarSyncRepo.GetByID(c.Request.Context(), syncID)
	if err != nil || sync == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar sync configuration not found"})
		return
	}

	// Verify ownership
	if sync.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this sync configuration"})
		return
	}

	if err := h.googleCalendarSyncRepo.Delete(c.Request.Context(), syncID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete calendar sync configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Calendar sync configuration deleted successfully"})
}

// SyncNow triggers an immediate sync for a specific calendar
func (h *GoogleCalendarSyncHandler) SyncNow(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	syncID := c.Param("id")
	if syncID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sync ID is required"})
		return
	}

	// Get sync configuration
	sync, err := h.googleCalendarSyncRepo.GetByID(c.Request.Context(), syncID)
	if err != nil || sync == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar sync configuration not found"})
		return
	}

	// Verify ownership
	if sync.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to sync this calendar"})
		return
	}

	// Get Google integration
	integration, err := h.googleIntegrationRepo.GetByID(c.Request.Context(), sync.GoogleIntegrationID)
	if err != nil || integration == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Google integration"})
		return
	}

	// Check if token needs refresh
	if integration.IsTokenExpiringSoon() {
		newTokens, err := h.oauth2Service.RefreshToken(c.Request.Context(), integration.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh Google token"})
			return
		}

		// Update tokens in database
		if err := h.googleIntegrationRepo.UpdateTokens(c.Request.Context(), integration.ID,
			newTokens.AccessToken, newTokens.RefreshToken, newTokens.Expiry); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		integration.AccessToken = newTokens.AccessToken
	}

	// Perform sync based on sync direction
	var syncError error
	syncedCount := 0

	switch sync.SyncDirection {
	case entities.SyncDirectionFromGoogle:
		syncedCount, syncError = h.syncFromGoogle(c.Request.Context(), sync, integration)
	case entities.SyncDirectionToGoogle:
		syncedCount, syncError = h.syncToGoogle(c.Request.Context(), sync, integration)
	case entities.SyncDirectionBidirectional:
		// First sync from Google, then to Google
		fromCount, fromErr := h.syncFromGoogle(c.Request.Context(), sync, integration)
		if fromErr == nil {
			toCount, toErr := h.syncToGoogle(c.Request.Context(), sync, integration)
			syncedCount = fromCount + toCount
			syncError = toErr
		} else {
			syncError = fromErr
			syncedCount = fromCount
		}
	}

	// Update sync status
	now := time.Now()
	var status entities.CalendarSyncStatus
	var errorMsg string

	if syncError != nil {
		status = entities.SyncStatusError
		errorMsg = syncError.Error()
	} else {
		status = entities.SyncStatusActive
		errorMsg = ""
	}

	if err := h.googleCalendarSyncRepo.UpdateSyncStatus(c.Request.Context(), syncID, status, &now, errorMsg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sync status"})
		return
	}

	if syncError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Sync completed with errors",
			"error_detail": syncError.Error(),
			"synced_count": syncedCount,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Sync completed successfully",
		"synced_count": syncedCount,
		"synced_at":    now,
	})
}

// syncFromGoogle syncs events from Google Calendar to local database
func (h *GoogleCalendarSyncHandler) syncFromGoogle(ctx context.Context, sync *entities.GoogleCalendarSync, integration *entities.GoogleIntegration) (int, error) {
	// Get events from Google Calendar
	now := time.Now()
	timeMin := now.AddDate(0, -1, 0) // 1 month ago
	timeMax := now.AddDate(0, 3, 0)  // 3 months ahead

	if !sync.Settings.SyncPastEvents {
		timeMin = now
	}

	googleEvents, err := h.calendarService.GetEvents(ctx, integration.AccessToken, sync.CalendarID, timeMin, timeMax)
	if err != nil {
		return 0, err
	}

	syncedCount := 0
	for _, googleEvent := range googleEvents {
		// Check if event already exists
		existingEvent, err := h.eventRepo.GetByExternalID(ctx, sync.UserID, googleEvent.ID)
		if err == nil && existingEvent != nil {
			// Update existing event
			existingEvent.Title = googleEvent.Summary
			existingEvent.Description = googleEvent.Description
			existingEvent.Location = googleEvent.Location
			existingEvent.StartTime = googleEvent.StartTime
			existingEvent.EndTime = googleEvent.EndTime
			existingEvent.UpdatedAt = time.Now()

			if err := h.eventRepo.Update(ctx, existingEvent); err == nil {
				syncedCount++
			}
		} else {
			// Create new event
			event := &entities.Event{
				ID:             entities.EventID(uuid.New().String()),
				UserID:         sync.UserID,
				Title:          googleEvent.Summary,
				Description:    googleEvent.Description,
				Location:       googleEvent.Location,
				StartTime:      googleEvent.StartTime,
				EndTime:        googleEvent.EndTime,
				Timezone:       "UTC", // TODO: Get from Google event
				Status:         entities.EventStatusConfirmed,
				ExternalID:     googleEvent.ID,
				ExternalSource: "google",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			if err := h.eventRepo.Create(ctx, event); err == nil {
				syncedCount++
			}
		}
	}

	return syncedCount, nil
}

// syncToGoogle syncs events from local database to Google Calendar
func (h *GoogleCalendarSyncHandler) syncToGoogle(ctx context.Context, sync *entities.GoogleCalendarSync, integration *entities.GoogleIntegration) (int, error) {
	// Get local events that need to be synced
	now := time.Now()
	timeMin := now.AddDate(0, -1, 0) // 1 month ago
	timeMax := now.AddDate(0, 3, 0)  // 3 months ahead

	if !sync.Settings.SyncPastEvents {
		timeMin = now
	}

	localEvents, err := h.eventRepo.GetByTimeRange(ctx, sync.UserID, timeMin, timeMax)
	if err != nil {
		return 0, err
	}

	syncedCount := 0
	for _, event := range localEvents {
		// Skip events already synced to Google
		if event.ExternalID != "" && event.ExternalSource == "google" {
			continue
		}

		// Create event in Google Calendar
		googleEvent, err := h.calendarService.CreateEvent(ctx, integration.AccessToken, sync.CalendarID, event)
		if err == nil {
			// Update local event with Google ID
			event.ExternalID = googleEvent.ID
			event.ExternalSource = "google"
			if err := h.eventRepo.Update(ctx, event); err == nil {
				syncedCount++
			}
		}
	}

	return syncedCount, nil
}

func (h *GoogleCalendarSyncHandler) toSyncConfigResponse(sync *entities.GoogleCalendarSync) CalendarSyncConfigResponse {
	var lastSyncAt *time.Time
	if !sync.LastSyncAt.IsZero() {
		lastSyncAt = &sync.LastSyncAt
	}
	
	return CalendarSyncConfigResponse{
		ID:            sync.ID,
		CalendarID:    sync.CalendarID,
		CalendarName:  sync.CalendarName,
		SyncDirection: sync.SyncDirection,
		SyncStatus:    sync.SyncStatus,
		LastSyncAt:    lastSyncAt,
		LastSyncError: sync.LastSyncError,
		Settings:      sync.Settings,
		CreatedAt:     sync.CreatedAt,
		UpdatedAt:     sync.UpdatedAt,
	}
}
