// +build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

// These tests require a real PostgreSQL database running
// Run with: go test -tags=integration ./internal/adapters/postgres/...

func TestUserRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	pool, cleanup := SetupTestDB(t, DefaultTestDBConfig())
	defer cleanup()

	repo := NewUserRepository(pool)
	ctx := context.Background()

	t.Run("Create and Get User", func(t *testing.T) {
		user := &entities.User{
			ID:    "integration-test-user",
			Email: "integration@test.com",
			Name:  "Integration Test User",
			Profile: entities.UserProfile{
				FirstName: "Integration",
				LastName:  "Test",
				Timezone:  "UTC",
			},
			Settings: entities.UserSettings{
				Language:            "en",
				DateFormat:          "YYYY-MM-DD",
				TimeFormat:          "24h",
				WeekStartDay:        1,
				NotificationEnabled: true,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Test Create
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Test GetByID
		retrieved, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Email, retrieved.Email)
		assert.Equal(t, user.Name, retrieved.Name)
		assert.Equal(t, user.Profile.FirstName, retrieved.Profile.FirstName)
		assert.Equal(t, user.Settings.Language, retrieved.Settings.Language)

		// Test GetByEmail
		byEmail, err := repo.GetByEmail(ctx, user.Email)
		require.NoError(t, err)
		require.NotNil(t, byEmail)
		assert.Equal(t, user.ID, byEmail.ID)

		// Test ExistsByEmail
		exists, err := repo.ExistsByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.True(t, exists)

		// Test Update
		user.Name = "Updated Name"
		user.Profile.FirstName = "Updated"
		err = repo.Update(ctx, user)
		require.NoError(t, err)

		updated, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
		assert.Equal(t, "Updated", updated.Profile.FirstName)

		// Test Delete
		err = repo.Delete(ctx, user.ID)
		require.NoError(t, err)

		deleted, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Nil(t, deleted)
	})
}

func TestGoalRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	pool, cleanup := SetupTestDB(t, DefaultTestDBConfig())
	defer cleanup()

	// Create test user first
	userID := CreateTestUser(t, pool)

	repo := NewGoalRepository(pool)
	ctx := context.Background()

	t.Run("Create and Get Goal", func(t *testing.T) {
		deadline := time.Now().Add(30 * 24 * time.Hour)
		goal := &entities.Goal{
			ID:          "integration-test-goal",
			UserID:      entities.UserID(userID),
			Title:       "Integration Test Goal",
			Description: "This is an integration test goal",
			Category:    entities.GoalCategoryEducation,
			Priority:    entities.PriorityHigh,
			Status:      entities.GoalStatusActive,
			Progress:    25,
			Deadline:    &deadline,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Test Create
		err := repo.Create(ctx, goal)
		require.NoError(t, err)

		// Test GetByID
		retrieved, err := repo.GetByID(ctx, goal.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		assert.Equal(t, goal.ID, retrieved.ID)
		assert.Equal(t, goal.UserID, retrieved.UserID)
		assert.Equal(t, goal.Title, retrieved.Title)
		assert.Equal(t, goal.Category, retrieved.Category)
		assert.Equal(t, goal.Priority, retrieved.Priority)
		assert.Equal(t, goal.Status, retrieved.Status)
		assert.Equal(t, goal.Progress, retrieved.Progress)

		// Deadline comparison (allowing for small time differences due to DB precision)
		if goal.Deadline != nil && retrieved.Deadline != nil {
			assert.WithinDuration(t, *goal.Deadline, *retrieved.Deadline, time.Second)
		}
	})
}

func TestMoodRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	pool, cleanup := SetupTestDB(t, DefaultTestDBConfig())
	defer cleanup()

	// Create test user first
	userID := CreateTestUser(t, pool)

	repo := NewMoodRepository(pool)
	ctx := context.Background()

	t.Run("Create and Get Mood", func(t *testing.T) {
		moodDate := time.Date(2025, 7, 27, 0, 0, 0, 0, time.UTC)
		mood := &entities.Mood{
			ID:         "integration-test-mood",
			UserID:     entities.UserID(userID),
			Date:       moodDate,
			Level:      entities.MoodLevelGood,
			Notes:      "Integration test mood",
			Tags:       []entities.MoodTag{entities.MoodTagWork, entities.MoodTagProductivity},
			RecordedAt: time.Now(),
		}

		// Test Create
		err := repo.Create(ctx, mood)
		require.NoError(t, err)

		// Test GetByID
		retrieved, err := repo.GetByID(ctx, mood.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		assert.Equal(t, mood.ID, retrieved.ID)
		assert.Equal(t, mood.UserID, retrieved.UserID)
		assert.Equal(t, mood.Level, retrieved.Level)
		assert.Equal(t, mood.Notes, retrieved.Notes)
		
		// Date comparison (only date part, not time)
		assert.True(t, retrieved.IsSameDate(moodDate))

		// Tags comparison
		assert.Len(t, retrieved.Tags, len(mood.Tags))
		for _, tag := range mood.Tags {
			assert.True(t, retrieved.HasTag(tag))
		}
	})
}