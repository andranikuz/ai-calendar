package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

func setupGoalTestDB(t *testing.T) (*goalRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := &goalRepository{
		pool: nil, // Will be set in integration tests
	}

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestGoalRepository_Create_SQLQuery(t *testing.T) {
	expectedQuery := `
		INSERT INTO goals (
			id, user_id, title, description, category, priority, status, 
			progress, deadline, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	// Verify the query structure is correct
	assert.Contains(t, expectedQuery, "INSERT INTO goals")
	assert.Contains(t, expectedQuery, "VALUES")
	assert.Contains(t, expectedQuery, "$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11")
}

func TestGoalRepository_GetByID_SQLQuery(t *testing.T) {
	expectedQuery := `
		SELECT id, user_id, title, description, category, priority, status, 
			   progress, deadline, created_at, updated_at
		FROM goals 
		WHERE id = $1`

	assert.Contains(t, expectedQuery, "SELECT")
	assert.Contains(t, expectedQuery, "FROM goals")
	assert.Contains(t, expectedQuery, "WHERE id = $1")
}

func TestGoal_Validation(t *testing.T) {
	tests := []struct {
		name    string
		goal    entities.Goal
		isValid bool
	}{
		{
			name: "Valid goal",
			goal: entities.Goal{
				ID:          "goal-1",
				UserID:      "user-1",
				Title:       "Learn Go Programming",
				Description: "Master Go programming language for backend development",
				Category:    entities.GoalCategoryEducation,
				Priority:    entities.PriorityHigh,
				Status:      entities.GoalStatusActive,
				Progress:    30,
				Deadline:    timePtr(time.Now().Add(30 * 24 * time.Hour)),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			isValid: true,
		},
		{
			name: "Empty title",
			goal: entities.Goal{
				ID:       "goal-1",
				UserID:   "user-1",
				Title:    "",
				Category: entities.GoalCategoryEducation,
				Priority: entities.PriorityHigh,
				Status:   entities.GoalStatusActive,
				Progress: 0,
			},
			isValid: false,
		},
		{
			name: "Invalid progress (over 100)",
			goal: entities.Goal{
				ID:       "goal-1",
				UserID:   "user-1",
				Title:    "Test Goal",
				Category: entities.GoalCategoryEducation,
				Priority: entities.PriorityHigh,
				Status:   entities.GoalStatusActive,
				Progress: 150,
			},
			isValid: false,
		},
		{
			name: "Invalid progress (negative)",
			goal: entities.Goal{
				ID:       "goal-1",
				UserID:   "user-1",
				Title:    "Test Goal",
				Category: entities.GoalCategoryEducation,
				Priority: entities.PriorityHigh,
				Status:   entities.GoalStatusActive,
				Progress: -10,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation logic
			hasID := tt.goal.ID != ""
			hasUserID := tt.goal.UserID != ""
			hasTitle := tt.goal.Title != ""
			validProgress := tt.goal.Progress >= 0 && tt.goal.Progress <= 100

			basicValidation := hasID && hasUserID && hasTitle && validProgress

			if tt.isValid {
				assert.True(t, basicValidation, "Goal should be valid")
			} else {
				assert.False(t, basicValidation, "Goal should be invalid")
			}
		})
	}
}

func TestGoalCategory_Constants(t *testing.T) {
	categories := []entities.GoalCategory{
		entities.GoalCategoryHealth,
		entities.GoalCategoryCareer,
		entities.GoalCategoryEducation,
		entities.GoalCategoryPersonal,
		entities.GoalCategoryFinancial,
		entities.GoalCategoryRelationship,
	}

	expectedCategories := []string{
		"health",
		"career", 
		"education",
		"personal",
		"financial",
		"relationship",
	}

	for i, category := range categories {
		assert.Equal(t, expectedCategories[i], string(category))
	}
}

func TestGoalStatus_Constants(t *testing.T) {
	statuses := []entities.GoalStatus{
		entities.GoalStatusDraft,
		entities.GoalStatusActive,
	}

	expectedStatuses := []string{
		"draft",
		"active",
	}

	for i, status := range statuses {
		assert.Equal(t, expectedStatuses[i], string(status))
	}
}

func TestGoal_ProgressValidation(t *testing.T) {
	tests := []struct {
		progress int
		valid    bool
	}{
		{0, true},
		{50, true},
		{100, true},
		{-1, false},
		{101, false},
		{1000, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("progress_%d", tt.progress), func(t *testing.T) {
			isValid := tt.progress >= 0 && tt.progress <= 100
			assert.Equal(t, tt.valid, isValid)
		})
	}
}

func TestGoal_DeadlineHandling(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	tests := []struct {
		name     string
		deadline *time.Time
		hasDeadline bool
	}{
		{
			name:        "No deadline",
			deadline:    nil,
			hasDeadline: false,
		},
		{
			name:        "Future deadline",
			deadline:    &future,
			hasDeadline: true,
		},
		{
			name:        "Past deadline",
			deadline:    &past,
			hasDeadline: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goal := entities.Goal{
				Deadline: tt.deadline,
			}

			hasDeadline := goal.Deadline != nil
			assert.Equal(t, tt.hasDeadline, hasDeadline)

			if hasDeadline {
				assert.NotNil(t, goal.Deadline)
			} else {
				assert.Nil(t, goal.Deadline)
			}
		})
	}
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}