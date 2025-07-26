package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/auth"
	appHandlers "github.com/andranikuz/smart-goal-calendar/internal/application/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/application/commands"
	"github.com/andranikuz/smart-goal-calendar/internal/application/queries"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/services"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
)

type UserHTTPHandler struct {
	userHandler *appHandlers.UserHandler
	userService *services.UserService
	jwtService  *auth.JWTService
}

func NewUserHTTPHandler(
	userHandler *appHandlers.UserHandler,
	userService *services.UserService,
	jwtService *auth.JWTService,
) *UserHTTPHandler {
	return &UserHTTPHandler{
		userHandler: userHandler,
		userService: userService,
		jwtService:  jwtService,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email     string                  `json:"email" binding:"required,email"`
	Name      string                  `json:"name" binding:"required,min=2,max=100"`
	Password  string                  `json:"password" binding:"required,min=8"`
	Profile   *entities.UserProfile   `json:"profile,omitempty"`
	Settings  *entities.UserSettings  `json:"settings,omitempty"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Name     *string                 `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Profile  *entities.UserProfile   `json:"profile,omitempty"`
	Settings *entities.UserSettings  `json:"settings,omitempty"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID       entities.UserID       `json:"id"`
	Email    string                `json:"email"`
	Name     string                `json:"name"`
	Profile  entities.UserProfile  `json:"profile"`
	Settings entities.UserSettings `json:"settings"`
}

// Register handles user registration
func (h *UserHTTPHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Sanitize input data
	req.Email = h.userService.SanitizeEmail(req.Email)
	req.Name = h.userService.SanitizeName(req.Name)
	
	// Set default profile and settings if not provided
	profile := entities.UserProfile{}
	if req.Profile != nil {
		profile = *req.Profile
	} else {
		profile = h.userService.GetDefaultUserProfile()
	}
	
	settings := entities.UserSettings{}
	if req.Settings != nil {
		settings = *req.Settings
	} else {
		settings = h.userService.GetDefaultUserSettings()
	}
	
	// Create user command
	cmd := commands.CreateUserCommand{
		Email:    req.Email,
		Name:     req.Name,
		Profile:  profile,
		Settings: settings,
	}
	
	// Execute command
	result, err := h.userHandler.HandleCreateUser(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "registration_failed",
			"message": err.Error(),
		})
		return
	}
	
	// Get created user
	userQuery := queries.GetUserByIDQuery{UserID: result.UserID}
	userResult, err := h.userHandler.HandleGetUserByID(c.Request.Context(), userQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "user_retrieval_failed",
			"message": "Failed to retrieve created user",
		})
		return
	}
	
	// Generate JWT tokens
	tokenPair, err := h.jwtService.GenerateTokenPair(userResult.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_generation_failed",
			"message": "Failed to generate authentication tokens",
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    h.mapUserToResponse(userResult.User),
		"tokens":  tokenPair,
	})
}

// Login handles user authentication
func (h *UserHTTPHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Sanitize email
	req.Email = h.userService.SanitizeEmail(req.Email)
	
	// Get user by email
	userQuery := queries.GetUserByEmailQuery{Email: req.Email}
	userResult, err := h.userHandler.HandleGetUserByEmail(c.Request.Context(), userQuery)
	if err != nil || userResult.User == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_credentials",
			"message": "Invalid email or password",
		})
		return
	}
	
	// TODO: Implement password verification
	// For now, skip password verification (will be implemented with password hashing)
	
	// Generate JWT tokens
	tokenPair, err := h.jwtService.GenerateTokenPair(userResult.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_generation_failed",
			"message": "Failed to generate authentication tokens",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    h.mapUserToResponse(userResult.User),
		"tokens":  tokenPair,
	})
}

// GetProfile gets current user profile
func (h *UserHTTPHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	query := queries.GetUserProfileQuery{UserID: userID}
	result, err := h.userHandler.HandleGetUserProfile(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "user_not_found",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"user": UserResponse{
			ID:       result.UserID,
			Email:    result.Email,
			Name:     result.Name,
			Profile:  result.Profile,
			Settings: result.Settings,
		},
	})
}

// UpdateProfile updates current user profile
func (h *UserHTTPHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Sanitize name if provided
	if req.Name != nil {
		sanitized := h.userService.SanitizeName(*req.Name)
		req.Name = &sanitized
	}
	
	// Create update command
	cmd := commands.UpdateUserCommand{
		UserID:   userID,
		Name:     req.Name,
		Profile:  req.Profile,
		Settings: req.Settings,
	}
	
	// Execute command
	result, err := h.userHandler.HandleUpdateUser(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "update_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Profile updated successfully",
		"updated_at": result.UpdatedAt,
	})
}

// DeleteAccount deletes current user account
func (h *UserHTTPHandler) DeleteAccount(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}
	
	// Check if user can be deleted
	if err := h.userService.CanDeleteUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "deletion_not_allowed",
			"message": err.Error(),
		})
		return
	}
	
	// Create delete command
	cmd := commands.DeleteUserCommand{UserID: userID}
	
	// Execute command
	result, err := h.userHandler.HandleDeleteUser(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "deletion_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Account deleted successfully",
		"deleted_at": result.DeletedAt,
	})
}

// RefreshToken refreshes JWT tokens using refresh token
func (h *UserHTTPHandler) RefreshToken(c *gin.Context) {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}
	
	// Validate refresh token
	claims, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_refresh_token",
			"message": err.Error(),
		})
		return
	}
	
	// Get user
	userQuery := queries.GetUserByIDQuery{UserID: claims.UserID}
	userResult, err := h.userHandler.HandleGetUserByID(c.Request.Context(), userQuery)
	if err != nil || userResult.User == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User not found",
		})
		return
	}
	
	// Generate new token pair
	tokenPair, err := h.jwtService.RefreshTokenPair(req.RefreshToken, userResult.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_refresh_failed",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Tokens refreshed successfully",
		"tokens":  tokenPair,
	})
}

// mapUserToResponse converts User entity to API response format
func (h *UserHTTPHandler) mapUserToResponse(user *entities.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Profile:  user.Profile,
		Settings: user.Settings,
	}
}