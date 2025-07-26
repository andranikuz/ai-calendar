package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type MoodService struct{}

func NewMoodService() *MoodService {
	return &MoodService{}
}

func (s *MoodService) ValidateMood(mood *entities.Mood) error {
	if !mood.IsValid() {
		return fmt.Errorf("invalid mood data")
	}

	if mood.Date.After(time.Now().AddDate(0, 0, 1)) {
		return fmt.Errorf("mood date cannot be more than 1 day in the future")
	}

	if mood.Date.Before(time.Now().AddDate(-1, 0, 0)) {
		return fmt.Errorf("mood date cannot be more than 1 year in the past")
	}

	if len(mood.Notes) > 1000 {
		return fmt.Errorf("mood notes cannot exceed 1000 characters")
	}

	if len(mood.Tags) > 10 {
		return fmt.Errorf("mood cannot have more than 10 tags")
	}

	for _, tag := range mood.Tags {
		if !s.IsValidTag(tag) {
			return fmt.Errorf("invalid mood tag: %s", tag)
		}
	}

	return nil
}

func (s *MoodService) IsValidTag(tag entities.MoodTag) bool {
	validTags := []entities.MoodTag{
		entities.MoodTagWork,
		entities.MoodTagFamily,
		entities.MoodTagHealth,
		entities.MoodTagSocial,
		entities.MoodTagExercise,
		entities.MoodTagSleep,
		entities.MoodTagStress,
		entities.MoodTagProductivity,
		entities.MoodTagRelaxation,
		entities.MoodTagCreativity,
	}

	for _, validTag := range validTags {
		if tag == validTag {
			return true
		}
	}

	return false
}

func (s *MoodService) SanitizeMoodNotes(notes string) string {
	return strings.TrimSpace(notes)
}

func (s *MoodService) GetMoodAdvice(level entities.MoodLevel, tags []entities.MoodTag) string {
	advice := ""

	switch level {
	case entities.MoodLevelVeryBad:
		advice = "Consider reaching out to someone you trust or taking small steps to improve your day."
	case entities.MoodLevelBad:
		advice = "Try some self-care activities or reflect on what might help improve your mood."
	case entities.MoodLevelNeutral:
		advice = "A neutral day is perfectly fine. Consider what small thing might bring you joy."
	case entities.MoodLevelGood:
		advice = "Great to see you're feeling good! Try to maintain positive habits."
	case entities.MoodLevelVeryGood:
		advice = "Wonderful! Consider what contributed to this great mood to repeat it."
	}

	// Add tag-specific advice
	for _, tag := range tags {
		switch tag {
		case entities.MoodTagStress:
			advice += " Consider stress management techniques like deep breathing or meditation."
		case entities.MoodTagSleep:
			advice += " Good sleep hygiene can significantly impact your mood."
		case entities.MoodTagExercise:
			advice += " Physical activity can boost endorphins and improve mood."
		case entities.MoodTagWork:
			advice += " Work-life balance is important for overall well-being."
		}
	}

	return advice
}

func (s *MoodService) CalculateMoodTrend(moods []*entities.Mood) string {
	if len(moods) < 2 {
		return "insufficient_data"
	}

	// Sort by date (ascending)
	// Note: Assuming moods are already sorted by date DESC from repository
	recent := make([]*entities.Mood, len(moods))
	for i, mood := range moods {
		recent[len(moods)-1-i] = mood
	}

	if len(recent) < 7 {
		return "insufficient_data"
	}

	// Calculate average for first half vs second half of recent moods
	mid := len(recent) / 2
	firstHalfAvg := s.calculateAverageLevel(recent[:mid])
	secondHalfAvg := s.calculateAverageLevel(recent[mid:])

	diff := secondHalfAvg - firstHalfAvg

	if diff > 0.5 {
		return "improving"
	} else if diff < -0.5 {
		return "declining"
	} else {
		return "stable"
	}
}

func (s *MoodService) calculateAverageLevel(moods []*entities.Mood) float64 {
	if len(moods) == 0 {
		return 0
	}

	total := 0
	for _, mood := range moods {
		total += int(mood.Level)
	}

	return float64(total) / float64(len(moods))
}

func (s *MoodService) GetMoodInsights(moods []*entities.Mood) map[string]interface{} {
	if len(moods) == 0 {
		return map[string]interface{}{
			"message": "No mood data available",
		}
	}

	insights := map[string]interface{}{}

	// Calculate averages by day of week
	dayAverages := s.calculateDayOfWeekAverages(moods)
	insights["day_of_week_averages"] = dayAverages

	// Find best and worst days
	bestDay, worstDay := s.findBestAndWorstDays(dayAverages)
	insights["best_day_of_week"] = bestDay
	insights["worst_day_of_week"] = worstDay

	// Tag correlations
	tagCorrelations := s.calculateTagCorrelations(moods)
	insights["tag_correlations"] = tagCorrelations

	// Trend
	trend := s.CalculateMoodTrend(moods)
	insights["trend"] = trend

	return insights
}

func (s *MoodService) calculateDayOfWeekAverages(moods []*entities.Mood) map[string]float64 {
	dayTotals := make(map[time.Weekday][]int)

	for _, mood := range moods {
		day := mood.Date.Weekday()
		dayTotals[day] = append(dayTotals[day], int(mood.Level))
	}

	averages := make(map[string]float64)
	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	for day := time.Sunday; day <= time.Saturday; day++ {
		if levels, exists := dayTotals[day]; exists && len(levels) > 0 {
			total := 0
			for _, level := range levels {
				total += level
			}
			averages[dayNames[int(day)]] = float64(total) / float64(len(levels))
		}
	}

	return averages
}

func (s *MoodService) findBestAndWorstDays(dayAverages map[string]float64) (string, string) {
	var bestDay, worstDay string
	var bestAvg, worstAvg float64 = 0, 6

	for day, avg := range dayAverages {
		if avg > bestAvg {
			bestAvg = avg
			bestDay = day
		}
		if avg < worstAvg {
			worstAvg = avg
			worstDay = day
		}
	}

	return bestDay, worstDay
}

func (s *MoodService) calculateTagCorrelations(moods []*entities.Mood) map[string]float64 {
	tagTotals := make(map[entities.MoodTag][]int)

	for _, mood := range moods {
		for _, tag := range mood.Tags {
			tagTotals[tag] = append(tagTotals[tag], int(mood.Level))
		}
	}

	correlations := make(map[string]float64)

	for tag, levels := range tagTotals {
		if len(levels) > 0 {
			total := 0
			for _, level := range levels {
				total += level
			}
			correlations[string(tag)] = float64(total) / float64(len(levels))
		}
	}

	return correlations
}