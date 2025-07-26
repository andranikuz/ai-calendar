package repositories

import (
	"context"
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type MoodRepository interface {
	// Create a new mood entry
	Create(ctx context.Context, mood *entities.Mood) error
	
	// Get mood by ID
	GetByID(ctx context.Context, id entities.MoodID) (*entities.Mood, error)
	
	// Get mood for a specific date and user
	GetByUserIDAndDate(ctx context.Context, userID entities.UserID, date time.Time) (*entities.Mood, error)
	
	// Get all moods for a user
	GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Mood, error)
	
	// Get moods for a user within a date range
	GetByUserIDAndDateRange(ctx context.Context, userID entities.UserID, start, end time.Time) ([]*entities.Mood, error)
	
	// Get moods by level for a user
	GetByUserIDAndLevel(ctx context.Context, userID entities.UserID, level entities.MoodLevel) ([]*entities.Mood, error)
	
	// Get moods with specific tags
	GetByUserIDAndTags(ctx context.Context, userID entities.UserID, tags []entities.MoodTag) ([]*entities.Mood, error)
	
	// Get latest mood entry for a user
	GetLatestByUserID(ctx context.Context, userID entities.UserID) (*entities.Mood, error)
	
	// Update mood entry
	Update(ctx context.Context, mood *entities.Mood) error
	
	// Delete mood entry
	Delete(ctx context.Context, id entities.MoodID) error
	
	// Upsert mood (create or update for a specific date)
	UpsertByDate(ctx context.Context, mood *entities.Mood) error
	
	// Get mood statistics for a user
	GetStatsByUserID(ctx context.Context, userID entities.UserID, start, end time.Time) (*MoodStats, error)
	
	// Get mood trends for a user
	GetTrendsByUserID(ctx context.Context, userID entities.UserID, days int) ([]*MoodTrend, error)
	
	// Check if mood exists for a specific date
	ExistsByUserIDAndDate(ctx context.Context, userID entities.UserID, date time.Time) (bool, error)
}

// MoodStats represents aggregated mood statistics
type MoodStats struct {
	UserID         entities.UserID `json:"user_id"`
	TotalEntries   int             `json:"total_entries"`
	AverageLevel   float64         `json:"average_level"`
	MostCommonTag  entities.MoodTag `json:"most_common_tag"`
	BestDay        time.Time       `json:"best_day"`
	WorstDay       time.Time       `json:"worst_day"`
	LevelCounts    map[entities.MoodLevel]int `json:"level_counts"`
	TagCounts      map[entities.MoodTag]int   `json:"tag_counts"`
}

// MoodTrend represents mood data for trend analysis
type MoodTrend struct {
	Date   time.Time           `json:"date"`
	Level  entities.MoodLevel  `json:"level"`
	Tags   []entities.MoodTag  `json:"tags"`
	Notes  string              `json:"notes"`
}