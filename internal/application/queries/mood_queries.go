package queries

import (
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type GetMoodByIDQuery struct {
	ID entities.MoodID `json:"id"`
}

type GetMoodByUserIDAndDateQuery struct {
	UserID entities.UserID `json:"user_id"`
	Date   time.Time       `json:"date"`
}

type GetMoodsByUserIDQuery struct {
	UserID entities.UserID `json:"user_id"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

type GetMoodsByDateRangeQuery struct {
	UserID entities.UserID `json:"user_id"`
	Start  time.Time       `json:"start"`
	End    time.Time       `json:"end"`
}

type GetMoodsByLevelQuery struct {
	UserID entities.UserID    `json:"user_id"`
	Level  entities.MoodLevel `json:"level"`
}

type GetMoodsByTagsQuery struct {
	UserID entities.UserID     `json:"user_id"`
	Tags   []entities.MoodTag  `json:"tags"`
}

type GetLatestMoodQuery struct {
	UserID entities.UserID `json:"user_id"`
}

type GetMoodStatsQuery struct {
	UserID entities.UserID `json:"user_id"`
	Start  time.Time       `json:"start"`
	End    time.Time       `json:"end"`
}

type GetMoodTrendsQuery struct {
	UserID entities.UserID `json:"user_id"`
	Days   int             `json:"days"`
}

type CheckMoodExistsQuery struct {
	UserID entities.UserID `json:"user_id"`
	Date   time.Time       `json:"date"`
}