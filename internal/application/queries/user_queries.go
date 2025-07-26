package queries

import (
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

// GetUserByIDQuery represents a query to get user by ID
type GetUserByIDQuery struct {
	UserID entities.UserID `json:"user_id" validate:"required"`
}

// GetUserByEmailQuery represents a query to get user by email
type GetUserByEmailQuery struct {
	Email string `json:"email" validate:"required,email"`
}

// GetUserProfileQuery represents a query to get user profile
type GetUserProfileQuery struct {
	UserID entities.UserID `json:"user_id" validate:"required"`
}

// Results
type GetUserResult struct {
	User *entities.User `json:"user"`
}

type GetUserProfileResult struct {
	UserID   entities.UserID       `json:"user_id"`
	Email    string                `json:"email"`
	Name     string                `json:"name"`
	Profile  entities.UserProfile  `json:"profile"`
	Settings entities.UserSettings `json:"settings"`
}