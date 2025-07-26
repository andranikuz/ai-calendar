package commands

import (
	"time"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type CreateMoodCommand struct {
	UserID entities.UserID       `json:"user_id"`
	Date   time.Time             `json:"date"`
	Level  entities.MoodLevel    `json:"level"`
	Notes  string                `json:"notes"`
	Tags   []entities.MoodTag    `json:"tags"`
}

type UpdateMoodCommand struct {
	ID     entities.MoodID       `json:"id"`
	UserID entities.UserID       `json:"user_id"`
	Level  entities.MoodLevel    `json:"level"`
	Notes  string                `json:"notes"`
	Tags   []entities.MoodTag    `json:"tags"`
}

type DeleteMoodCommand struct {
	ID     entities.MoodID       `json:"id"`
	UserID entities.UserID       `json:"user_id"`
}

type UpsertMoodByDateCommand struct {
	UserID entities.UserID       `json:"user_id"`
	Date   time.Time             `json:"date"`
	Level  entities.MoodLevel    `json:"level"`
	Notes  string                `json:"notes"`
	Tags   []entities.MoodTag    `json:"tags"`
}