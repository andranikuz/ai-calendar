package repositories

import (
	"context"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type UserRepository interface {
	// Create a new user
	Create(ctx context.Context, user *entities.User) error
	
	// Get user by ID
	GetByID(ctx context.Context, id entities.UserID) (*entities.User, error)
	
	// Get user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	
	// Update user
	Update(ctx context.Context, user *entities.User) error
	
	// Delete user
	Delete(ctx context.Context, id entities.UserID) error
	
	// Check if user exists by email
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}