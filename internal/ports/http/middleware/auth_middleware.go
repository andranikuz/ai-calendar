package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/auth"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type AuthMiddleware struct {
	jwtService *auth.JWTService
}

func NewAuthMiddleware(jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth middleware that requires valid JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "authorization_required",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}
		
		// Extract token from Bearer header
		token, err := m.jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_authorization_header",
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}
		
		// Validate token
		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}
		
		// Check if token is expired
		if m.jwtService.IsTokenExpired(claims) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "token_expired",
				"message": "Token has expired",
			})
			c.Abort()
			return
		}
		
		// Store user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("token_claims", claims)
		
		c.Next()
	}
}

// OptionalAuth middleware that extracts user info if token is present but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		
		token, err := m.jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}
		
		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.Next()
			return
		}
		
		if !m.jwtService.IsTokenExpired(claims) {
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_name", claims.Name)
			c.Set("token_claims", claims)
		}
		
		c.Next()
	}
}

// GetCurrentUserID extracts current user ID from gin context
func GetCurrentUserID(c *gin.Context) (entities.UserID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	
	id, ok := userID.(entities.UserID)
	return id, ok
}

// GetCurrentUserEmail extracts current user email from gin context
func GetCurrentUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}
	
	emailStr, ok := email.(string)
	return emailStr, ok
}

// GetCurrentUserName extracts current user name from gin context
func GetCurrentUserName(c *gin.Context) (string, bool) {
	name, exists := c.Get("user_name")
	if !exists {
		return "", false
	}
	
	nameStr, ok := name.(string)
	return nameStr, ok
}

// GetTokenClaims extracts JWT claims from gin context
func GetTokenClaims(c *gin.Context) (*auth.JWTClaims, bool) {
	claims, exists := c.Get("token_claims")
	if !exists {
		return nil, false
	}
	
	jwtClaims, ok := claims.(*auth.JWTClaims)
	return jwtClaims, ok
}

// RequireUserOwnership middleware that ensures the user can only access their own resources
func (m *AuthMiddleware) RequireUserOwnership(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user ID from token
		currentUserID, exists := GetCurrentUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}
		
		// Get requested user ID from URL parameter
		requestedUserID := entities.UserID(c.Param(paramName))
		if requestedUserID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "missing_parameter",
				"message": "Missing user ID parameter",
			})
			c.Abort()
			return
		}
		
		// Check if user is accessing their own resource
		if currentUserID != requestedUserID {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "access_denied",
				"message": "You can only access your own resources",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Allow specific origins in production
		allowedOrigins := []string{
			"http://localhost:3000",  // React dev server
			"http://localhost:8080",  // API server
			"https://smart-goal-calendar.com", // Production domain
		}
		
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}
		
		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// RequestLogger middleware for structured logging
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format can be implemented here
		return ""
	})
}

// RateLimiter middleware (basic implementation)
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Basic rate limiting implementation
		// In production, use a proper rate limiter like go-redis/redis_rate
		
		// Skip rate limiting for health checks
		if strings.HasPrefix(c.Request.URL.Path, "/health") {
			c.Next()
			return
		}
		
		// TODO: Implement proper rate limiting logic
		// For now, just pass through
		c.Next()
	}
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		c.Next()
	}
}