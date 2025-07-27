package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/andranikuz/smart-goal-calendar/internal/adapters/migrations"
)

// TestDBConfig holds test database configuration
type TestDBConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// DefaultTestDBConfig returns default test database configuration
func DefaultTestDBConfig() TestDBConfig {
	return TestDBConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "smart_goal_calendar_test",
		Username: "postgres",
		Password: "postgres",
	}
}

// SetupTestDB creates a test database connection and runs migrations
// This is used for integration tests that need a real database
func SetupTestDB(t *testing.T, config TestDBConfig) (*pgxpool.Pool, func()) {
	t.Helper()

	// Skip integration tests if not running with -integration flag
	if !testing.Short() {
		t.Skip("Skipping integration test")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err, "Failed to connect to test database")

	// Ping to ensure connection is working
	err = pool.Ping(ctx)
	require.NoError(t, err, "Failed to ping test database")

	// Run migrations to ensure schema is up to date
	migrator := migrations.NewMigrator(pool)
	err = migrator.Migrate()
	require.NoError(t, err, "Failed to run migrations on test database")

	cleanup := func() {
		// Clean up test data
		cleanupTestData(t, pool)
		pool.Close()
	}

	return pool, cleanup
}

// cleanupTestData removes all test data from the database
func cleanupTestData(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()
	
	// Clean up in reverse order due to foreign key constraints
	tables := []string{
		"google_calendar_syncs",
		"google_integrations", 
		"moods",
		"events",
		"milestones",
		"tasks",
		"goals",
		"users",
	}

	for _, table := range tables {
		_, err := pool.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Logf("Warning: Failed to clean up table %s: %v", table, err)
		}
	}
}

// CreateTestUser creates a test user for integration tests
func CreateTestUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()

	userID := "test-user-" + t.Name()
	ctx := context.Background()

	query := `
		INSERT INTO users (id, email, name, profile, settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`

	profile := `{"first_name":"Test","last_name":"User","timezone":"UTC"}`
	settings := `{"language":"en","date_format":"YYYY-MM-DD","time_format":"24h","week_start_day":1,"notification_enabled":true}`

	_, err := pool.Exec(ctx, query,
		userID,
		"test@example.com",
		"Test User",
		profile,
		settings,
	)
	require.NoError(t, err, "Failed to create test user")

	return userID
}

// CreateTestGoal creates a test goal for integration tests
func CreateTestGoal(t *testing.T, pool *pgxpool.Pool, userID string) string {
	t.Helper()

	goalID := "test-goal-" + t.Name()
	ctx := context.Background()

	query := `
		INSERT INTO goals (id, user_id, title, description, category, priority, status, progress, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`

	_, err := pool.Exec(ctx, query,
		goalID,
		userID,
		"Test Goal",
		"This is a test goal",
		"education",
		"high",
		"active",
		0,
	)
	require.NoError(t, err, "Failed to create test goal")

	return goalID
}

// CreateTestEvent creates a test event for integration tests
func CreateTestEvent(t *testing.T, pool *pgxpool.Pool, userID string) string {
	t.Helper()

	eventID := "test-event-" + t.Name()
	ctx := context.Background()

	query := `
		INSERT INTO events (id, user_id, title, description, start_time, end_time, timezone, location, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`

	_, err := pool.Exec(ctx, query,
		eventID,
		userID,
		"Test Event",
		"This is a test event",
		"2025-07-27 10:00:00",
		"2025-07-27 11:00:00",
		"UTC",
		"Test Location",
		"confirmed",
	)
	require.NoError(t, err, "Failed to create test event")

	return eventID
}

// CreateTestMood creates a test mood for integration tests
func CreateTestMood(t *testing.T, pool *pgxpool.Pool, userID string) string {
	t.Helper()

	moodID := "test-mood-" + t.Name()
	ctx := context.Background()

	query := `
		INSERT INTO moods (id, user_id, date, level, notes, tags, recorded_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())`

	_, err := pool.Exec(ctx, query,
		moodID,
		userID,
		"2025-07-27",
		4, // Good mood
		"Test mood notes",
		`["work", "productivity"]`,
	)
	require.NoError(t, err, "Failed to create test mood")

	return moodID
}