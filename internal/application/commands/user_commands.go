package commands

import (
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

// CreateUserCommand represents a command to create a new user
type CreateUserCommand struct {
	Email     string                   `json:"email" validate:"required,email"`
	Name      string                   `json:"name" validate:"required,min=2,max=100"`
	Profile   entities.UserProfile     `json:"profile"`
	Settings  entities.UserSettings    `json:"settings"`
}

// UpdateUserCommand represents a command to update user information
type UpdateUserCommand struct {
	UserID    entities.UserID          `json:"user_id" validate:"required"`
	Email     *string                  `json:"email,omitempty" validate:"omitempty,email"`
	Name      *string                  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Profile   *entities.UserProfile    `json:"profile,omitempty"`
	Settings  *entities.UserSettings   `json:"settings,omitempty"`
}

// DeleteUserCommand represents a command to delete a user
type DeleteUserCommand struct {
	UserID entities.UserID `json:"user_id" validate:"required"`
}

// CreateUserResult represents the result of creating a user
type CreateUserResult struct {
	UserID    entities.UserID `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
}

// UpdateUserResult represents the result of updating a user
type UpdateUserResult struct {
	UpdatedAt time.Time `json:"updated_at"`
}

// DeleteUserResult represents the result of deleting a user
type DeleteUserResult struct {
	DeletedAt time.Time `json:"deleted_at"`
}