package entities

import (
	"time"
)

type MoodID string

type Mood struct {
	ID         MoodID    `json:"id"`
	UserID     UserID    `json:"user_id"`
	Date       time.Time `json:"date"` // Date without time (YYYY-MM-DD)
	Level      MoodLevel `json:"level"`
	Notes      string    `json:"notes"`
	Tags       []MoodTag `json:"tags"`
	RecordedAt time.Time `json:"recorded_at"`
}

type MoodLevel int

const (
	MoodLevelVeryBad MoodLevel = 1
	MoodLevelBad     MoodLevel = 2
	MoodLevelNeutral MoodLevel = 3
	MoodLevelGood    MoodLevel = 4
	MoodLevelVeryGood MoodLevel = 5
)

func (ml MoodLevel) String() string {
	switch ml {
	case MoodLevelVeryBad:
		return "very_bad"
	case MoodLevelBad:
		return "bad"
	case MoodLevelNeutral:
		return "neutral"
	case MoodLevelGood:
		return "good"
	case MoodLevelVeryGood:
		return "very_good"
	default:
		return "unknown"
	}
}

func (ml MoodLevel) Emoji() string {
	switch ml {
	case MoodLevelVeryBad:
		return "😢"
	case MoodLevelBad:
		return "🙁"
	case MoodLevelNeutral:
		return "😐"
	case MoodLevelGood:
		return "🙂"
	case MoodLevelVeryGood:
		return "😄"
	default:
		return "❓"
	}
}

func (ml MoodLevel) IsValid() bool {
	return ml >= MoodLevelVeryBad && ml <= MoodLevelVeryGood
}

type MoodTag string

const (
	MoodTagWork        MoodTag = "work"
	MoodTagFamily      MoodTag = "family"
	MoodTagHealth      MoodTag = "health"
	MoodTagSocial      MoodTag = "social"
	MoodTagExercise    MoodTag = "exercise"
	MoodTagSleep       MoodTag = "sleep"
	MoodTagStress      MoodTag = "stress"
	MoodTagProductivity MoodTag = "productivity"
	MoodTagRelaxation  MoodTag = "relaxation"
	MoodTagCreativity  MoodTag = "creativity"
)

// Validation methods
func (m *Mood) IsValid() bool {
	return m.UserID != "" && 
		   !m.Date.IsZero() && 
		   m.Level.IsValid()
}

func (m *Mood) IsSameDate(date time.Time) bool {
	y1, m1, d1 := m.Date.Date()
	y2, m2, d2 := date.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (m *Mood) HasTag(tag MoodTag) bool {
	for _, t := range m.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (m *Mood) AddTag(tag MoodTag) {
	if !m.HasTag(tag) {
		m.Tags = append(m.Tags, tag)
	}
}

func (m *Mood) RemoveTag(tag MoodTag) {
	for i, t := range m.Tags {
		if t == tag {
			m.Tags = append(m.Tags[:i], m.Tags[i+1:]...)
			break
		}
	}
}