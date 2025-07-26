package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type JWTService struct {
	secretKey       []byte
	accessExpiry    time.Duration
	refreshExpiry   time.Duration
	issuer          string
}

type JWTClaims struct {
	UserID entities.UserID `json:"user_id"`
	Email  string          `json:"email"`
	Name   string          `json:"name"`
	Type   string          `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration, issuer string) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        issuer,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (s *JWTService) GenerateTokenPair(user *entities.User) (*TokenPair, error) {
	now := time.Now()
	
	// Generate access token
	accessClaims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   string(user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}
	
	// Generate refresh token
	refreshClaims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   string(user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}
	
	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    now.Add(s.accessExpiry),
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken validates and parses a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	
	// Additional validation
	if claims.Issuer != s.issuer {
		return nil, fmt.Errorf("invalid token issuer")
	}
	
	return claims, nil
}

// ValidateAccessToken validates specifically an access token
func (s *JWTService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	if claims.Type != "access" {
		return nil, fmt.Errorf("not an access token")
	}
	
	return claims, nil
}

// ValidateRefreshToken validates specifically a refresh token
func (s *JWTService) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	if claims.Type != "refresh" {
		return nil, fmt.Errorf("not a refresh token")
	}
	
	return claims, nil
}

// RefreshTokenPair generates new token pair using refresh token
func (s *JWTService) RefreshTokenPair(refreshTokenString string, user *entities.User) (*TokenPair, error) {
	// Validate refresh token
	_, err := s.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	
	// Generate new token pair
	return s.GenerateTokenPair(user)
}

// ExtractTokenFromHeader extracts token from Authorization header
func (s *JWTService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}
	
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf("invalid authorization header format")
	}
	
	return authHeader[len(bearerPrefix):], nil
}

// GetTokenExpiryTime returns the expiry time for access tokens
func (s *JWTService) GetTokenExpiryTime() time.Duration {
	return s.accessExpiry
}

// GetRefreshTokenExpiryTime returns the expiry time for refresh tokens
func (s *JWTService) GetRefreshTokenExpiryTime() time.Duration {
	return s.refreshExpiry
}

// IsTokenExpired checks if token is expired
func (s *JWTService) IsTokenExpired(claims *JWTClaims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}

// GetUserIDFromToken extracts user ID from token claims
func (s *JWTService) GetUserIDFromToken(tokenString string) (entities.UserID, error) {
	claims, err := s.ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}
	
	return claims.UserID, nil
}