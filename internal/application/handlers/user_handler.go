package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/andranikuz/smart-goal-calendar/internal/application/commands"
	"github.com/andranikuz/smart-goal-calendar/internal/application/queries"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
)

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// Command Handlers

func (h *UserHandler) HandleCreateUser(ctx context.Context, cmd commands.CreateUserCommand) (*commands.CreateUserResult, error) {
	// Check if user already exists
	exists, err := h.userRepo.ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", cmd.Email)
	}
	
	// Create new user entity
	now := time.Now()
	user := &entities.User{
		ID:        entities.UserID(uuid.New().String()),
		Email:     cmd.Email,
		Name:      cmd.Name,
		Profile:   cmd.Profile,
		Settings:  cmd.Settings,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	// Set default settings if not provided
	if user.Settings.Language == "" {
		user.Settings.Language = "en"
	}
	if user.Settings.DateFormat == "" {
		user.Settings.DateFormat = "YYYY-MM-DD"
	}
	if user.Settings.TimeFormat == "" {
		user.Settings.TimeFormat = "24h"
	}
	if user.Settings.WeekStartDay == 0 {
		user.Settings.WeekStartDay = 1 // Monday
	}
	
	// Set default profile timezone if not provided
	if user.Profile.Timezone == "" {
		user.Profile.Timezone = "UTC"
	}
	
	// Save user to repository
	if err := h.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return &commands.CreateUserResult{
		UserID:    user.ID,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (h *UserHandler) HandleUpdateUser(ctx context.Context, cmd commands.UpdateUserCommand) (*commands.UpdateUserResult, error) {
	// Get existing user
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	
	// Update fields if provided
	if cmd.Email != nil {
		// Check if new email is already taken by another user
		existingUser, err := h.userRepo.GetByEmail(ctx, *cmd.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email availability: %w", err)
		}
		
		if existingUser != nil && existingUser.ID != cmd.UserID {
			return nil, fmt.Errorf("email %s is already taken", *cmd.Email)
		}
		
		user.Email = *cmd.Email
	}
	
	if cmd.Name != nil {
		user.Name = *cmd.Name
	}
	
	if cmd.Profile != nil {
		user.Profile = *cmd.Profile
	}
	
	if cmd.Settings != nil {
		user.Settings = *cmd.Settings
	}
	
	// Update timestamp
	user.UpdatedAt = time.Now()
	
	// Save updated user
	if err := h.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	
	return &commands.UpdateUserResult{
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (h *UserHandler) HandleDeleteUser(ctx context.Context, cmd commands.DeleteUserCommand) (*commands.DeleteUserResult, error) {
	// Check if user exists
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	
	// Delete user
	if err := h.userRepo.Delete(ctx, cmd.UserID); err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	
	return &commands.DeleteUserResult{
		DeletedAt: time.Now(),
	}, nil
}

// Query Handlers

func (h *UserHandler) HandleGetUserByID(ctx context.Context, query queries.GetUserByIDQuery) (*queries.GetUserResult, error) {
	user, err := h.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return &queries.GetUserResult{
		User: user,
	}, nil
}

func (h *UserHandler) HandleGetUserByEmail(ctx context.Context, query queries.GetUserByEmailQuery) (*queries.GetUserResult, error) {
	user, err := h.userRepo.GetByEmail(ctx, query.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return &queries.GetUserResult{
		User: user,
	}, nil
}

func (h *UserHandler) HandleGetUserProfile(ctx context.Context, query queries.GetUserProfileQuery) (*queries.GetUserProfileResult, error) {
	user, err := h.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	
	return &queries.GetUserProfileResult{
		UserID:   user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Profile:  user.Profile,
		Settings: user.Settings,
	}, nil
}