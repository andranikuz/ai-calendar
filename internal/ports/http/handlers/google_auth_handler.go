package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/andranikuz/smart-goal-calendar/internal/adapters/google"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

type GoogleAuthHandler struct {
	oauth2Service               *google.OAuth2Service
	calendarService             *google.CalendarService
	googleIntegrationRepo       repositories.GoogleIntegrationRepository
	googleCalendarSyncRepo      repositories.GoogleCalendarSyncRepository
}

func NewGoogleAuthHandler(
	oauth2Service *google.OAuth2Service,
	calendarService *google.CalendarService,
	googleIntegrationRepo repositories.GoogleIntegrationRepository,
	googleCalendarSyncRepo repositories.GoogleCalendarSyncRepository,
) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		oauth2Service:               oauth2Service,
		calendarService:             calendarService,
		googleIntegrationRepo:       googleIntegrationRepo,
		googleCalendarSyncRepo:      googleCalendarSyncRepo,
	}
}

type AuthURLResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

func (h *GoogleAuthHandler) GetAuthURL(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Generate a unique state parameter for CSRF protection
	state := uuid.New().String() + ":" + string(userID)
	
	authURL := h.oauth2Service.GetAuthURL(state)

	c.JSON(http.StatusOK, AuthURLResponse{
		AuthURL: authURL,
		State:   state,
	})
}

type OAuthCallbackRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

type IntegrationResponse struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Enabled      bool      `json:"enabled"`
	CalendarID   string    `json:"calendar_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (h *GoogleAuthHandler) HandleCallback(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req OAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate state parameter (basic CSRF protection)
	if req.State == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	// Exchange authorization code for tokens
	tokens, err := h.oauth2Service.ExchangeCode(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange authorization code"})
		return
	}

	// Get user info from Google
	userInfo, err := h.oauth2Service.GetUserInfo(c.Request.Context(), tokens.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info from Google"})
		return
	}

	// Get primary calendar
	calendars, err := h.calendarService.GetCalendars(c.Request.Context(), tokens.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendars"})
		return
	}

	var primaryCalendarID string
	for _, calendar := range calendars {
		if calendar.Primary {
			primaryCalendarID = calendar.ID
			break
		}
	}

	// Check if integration already exists
	existingIntegration, err := h.googleIntegrationRepo.GetByGoogleUserID(c.Request.Context(), userInfo.UserID)
	if err == nil && existingIntegration != nil {
		// Update existing integration
		existingIntegration.AccessToken = tokens.AccessToken
		existingIntegration.RefreshToken = tokens.RefreshToken
		existingIntegration.TokenType = tokens.TokenType
		existingIntegration.ExpiresAt = tokens.Expiry
		existingIntegration.Email = userInfo.Email
		existingIntegration.Name = userInfo.Name
		existingIntegration.CalendarID = primaryCalendarID
		existingIntegration.Enabled = true
		existingIntegration.UpdatedAt = time.Now()

		if err := h.googleIntegrationRepo.Update(c.Request.Context(), existingIntegration); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Google integration"})
			return
		}

		response := h.toIntegrationResponse(existingIntegration)
		c.JSON(http.StatusOK, response)
		return
	}

	// Create new integration
	now := time.Now()
	integration := &entities.GoogleIntegration{
		ID:           entities.GoogleIntegrationID(uuid.New().String()),
		UserID:       userID,
		GoogleUserID: userInfo.UserID,
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresAt:    tokens.Expiry,
		Scopes:       []string{"calendar.readonly", "calendar.events"},
		CalendarID:   primaryCalendarID,
		Enabled:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.googleIntegrationRepo.Create(c.Request.Context(), integration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Google integration"})
		return
	}

	// Create default calendar sync configuration
	if primaryCalendarID != "" {
		sync := &entities.GoogleCalendarSync{
			ID:                  uuid.New().String(),
			UserID:              userID,
			GoogleIntegrationID: integration.ID,
			CalendarID:          primaryCalendarID,
			CalendarName:        "Primary",
			SyncDirection:       entities.SyncDirectionBidirectional,
			SyncStatus:          entities.SyncStatusActive,
			Settings:            entities.DefaultCalendarSyncSettings(),
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		if err := h.googleCalendarSyncRepo.Create(c.Request.Context(), sync); err != nil {
			// Log error but don't fail the integration creation
			// TODO: Add proper logging
		}
	}

	response := h.toIntegrationResponse(integration)
	c.JSON(http.StatusCreated, response)
}

func (h *GoogleAuthHandler) GetIntegration(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	integration, err := h.googleIntegrationRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Google integration"})
		return
	}

	if integration == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Google integration not found"})
		return
	}

	response := h.toIntegrationResponse(integration)
	c.JSON(http.StatusOK, response)
}

func (h *GoogleAuthHandler) DisconnectIntegration(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	integration, err := h.googleIntegrationRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Google integration"})
		return
	}

	if integration == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Google integration not found"})
		return
	}

	// Delete calendar sync configurations
	syncs, err := h.googleCalendarSyncRepo.GetByIntegrationID(c.Request.Context(), integration.ID)
	if err == nil {
		for _, sync := range syncs {
			h.googleCalendarSyncRepo.Delete(c.Request.Context(), sync.ID)
		}
	}

	// Delete integration
	if err := h.googleIntegrationRepo.Delete(c.Request.Context(), integration.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Google integration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Google integration disconnected successfully"})
}

func (h *GoogleAuthHandler) GetCalendars(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	integration, err := h.googleIntegrationRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Google integration"})
		return
	}

	if integration == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Google integration not found"})
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

	calendars, err := h.calendarService.GetCalendars(c.Request.Context(), integration.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendars from Google"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"calendars": calendars})
}

func (h *GoogleAuthHandler) toIntegrationResponse(integration *entities.GoogleIntegration) IntegrationResponse {
	return IntegrationResponse{
		ID:         string(integration.ID),
		Email:      integration.Email,
		Name:       integration.Name,
		Enabled:    integration.Enabled,
		CalendarID: integration.CalendarID,
		CreatedAt:  integration.CreatedAt,
	}
}