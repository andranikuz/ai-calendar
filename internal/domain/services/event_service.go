package services

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type EventService struct{}

func NewEventService() *EventService {
	return &EventService{}
}

// SanitizeEventTitle cleans and validates event title
func (s *EventService) SanitizeEventTitle(title string) string {
	// Trim whitespace
	title = strings.TrimSpace(title)
	
	// Remove any control characters
	var cleaned strings.Builder
	for _, r := range title {
		if r != '\t' && r != '\n' && r != '\r' {
			cleaned.WriteRune(r)
		}
	}
	
	return cleaned.String()
}

// SanitizeEventDescription cleans and validates event description
func (s *EventService) SanitizeEventDescription(description string) string {
	// Trim whitespace
	description = strings.TrimSpace(description)
	
	// Allow newlines in description but remove other control characters
	var cleaned strings.Builder
	for _, r := range description {
		if r == '\n' || r == '\r' || (r >= 32 && r != 127) {
			cleaned.WriteRune(r)
		}
	}
	
	return cleaned.String()
}

// ValidateEventCreation validates event data for creation
func (s *EventService) ValidateEventCreation(event *entities.Event) error {
	// Validate title
	if event.Title == "" {
		return fmt.Errorf("event title is required")
	}
	
	if utf8.RuneCountInString(event.Title) > 255 {
		return fmt.Errorf("event title must be 255 characters or less")
	}
	
	// Validate description length
	if utf8.RuneCountInString(event.Description) > 1000 {
		return fmt.Errorf("event description must be 1000 characters or less")
	}
	
	// Validate time range
	if event.EndTime.Before(event.StartTime) || event.EndTime.Equal(event.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}
	
	// Validate event duration (max 24 hours)
	duration := event.EndTime.Sub(event.StartTime)
	if duration > 24*time.Hour {
		return fmt.Errorf("event duration cannot exceed 24 hours")
	}
	
	// Validate timezone
	if event.Timezone == "" {
		event.Timezone = "UTC" // Default to UTC
	}
	
	// Validate timezone format
	_, err := time.LoadLocation(event.Timezone)
	if err != nil {
		return fmt.Errorf("invalid timezone: %s", event.Timezone)
	}
	
	// Validate location length
	if utf8.RuneCountInString(event.Location) > 255 {
		return fmt.Errorf("event location must be 255 characters or less")
	}
	
	// Validate status
	if !s.IsValidEventStatus(event.Status) {
		return fmt.Errorf("invalid event status: %s", event.Status)
	}
	
	// Validate external ID length
	if utf8.RuneCountInString(event.ExternalID) > 255 {
		return fmt.Errorf("external ID must be 255 characters or less")
	}
	
	// Validate external source length
	if utf8.RuneCountInString(event.ExternalSource) > 50 {
		return fmt.Errorf("external source must be 50 characters or less")
	}
	
	// Validate attendees
	if len(event.Attendees) > 100 {
		return fmt.Errorf("event cannot have more than 100 attendees")
	}
	
	for _, attendee := range event.Attendees {
		if err := s.validateAttendee(attendee); err != nil {
			return fmt.Errorf("invalid attendee: %w", err)
		}
	}
	
	return nil
}

// IsValidEventStatus checks if the event status is valid
func (s *EventService) IsValidEventStatus(status entities.EventStatus) bool {
	switch status {
	case entities.EventStatusTentative, entities.EventStatusConfirmed, entities.EventStatusCancelled:
		return true
	default:
		return false
	}
}

// validateAttendee validates an attendee object
func (s *EventService) validateAttendee(attendee entities.Attendee) error {
	if attendee.Email == "" {
		return fmt.Errorf("attendee email is required")
	}
	
	if utf8.RuneCountInString(attendee.Email) > 255 {
		return fmt.Errorf("attendee email must be 255 characters or less")
	}
	
	if utf8.RuneCountInString(attendee.Name) > 255 {
		return fmt.Errorf("attendee name must be 255 characters or less")
	}
	
	// Basic email validation
	if !strings.Contains(attendee.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// CalculateEventDuration calculates the duration of an event in minutes
func (s *EventService) CalculateEventDuration(startTime, endTime time.Time) int {
	duration := endTime.Sub(startTime)
	return int(duration.Minutes())
}

// IsEventActive checks if an event is currently active
func (s *EventService) IsEventActive(event *entities.Event) bool {
	now := time.Now()
	return event.Status == entities.EventStatusConfirmed &&
		now.After(event.StartTime) &&
		now.Before(event.EndTime)
}

// IsEventUpcoming checks if an event is upcoming (starts within next 24 hours)
func (s *EventService) IsEventUpcoming(event *entities.Event) bool {
	now := time.Now()
	return event.Status != entities.EventStatusCancelled &&
		event.StartTime.After(now) &&
		event.StartTime.Before(now.Add(24*time.Hour))
}

// GetEventColor returns a color code based on event properties
func (s *EventService) GetEventColor(event *entities.Event) string {
	// Color by status
	switch event.Status {
	case entities.EventStatusTentative:
		return "#FFA500" // Orange
	case entities.EventStatusConfirmed:
		if event.GoalID != nil {
			return "#4CAF50" // Green for goal-linked events
		}
		return "#2196F3" // Blue for regular events
	case entities.EventStatusCancelled:
		return "#9E9E9E" // Gray
	default:
		return "#2196F3" // Default blue
	}
}

// FormatEventTimeRange formats event time range for display
func (s *EventService) FormatEventTimeRange(event *entities.Event, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC // Fallback to UTC
	}
	
	start := event.StartTime.In(loc)
	end := event.EndTime.In(loc)
	
	// Same day
	if start.Format("2006-01-02") == end.Format("2006-01-02") {
		return fmt.Sprintf("%s - %s",
			start.Format("Jan 2, 15:04"),
			end.Format("15:04")), nil
	}
	
	// Different days
	return fmt.Sprintf("%s - %s",
		start.Format("Jan 2, 15:04"),
		end.Format("Jan 2, 15:04")), nil
}

// GenerateEventSlug creates a URL-friendly slug from event title
func (s *EventService) GenerateEventSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace spaces and special characters with hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}
	
	// Clean up multiple hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	
	// Limit length
	if len(slug) > 50 {
		slug = slug[:50]
	}
	
	return slug
}

// CheckTimeSlotAvailability checks if a time slot is available for a user
func (s *EventService) CheckTimeSlotAvailability(startTime, endTime time.Time) error {
	// Validate time range
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return fmt.Errorf("end time must be after start time")
	}
	
	// Check if the time slot is in the past
	if startTime.Before(time.Now()) {
		return fmt.Errorf("cannot create events in the past")
	}
	
	// Check business hours (optional constraint)
	start := startTime.Hour()
	end := endTime.Hour()
	
	// Allow events outside business hours but warn about very early/late times
	if start < 6 || end > 23 {
		// This could be a warning rather than an error in a real application
		// For now, we'll allow it
	}
	
	return nil
}

// SuggestEventTimes suggests optimal times for an event based on duration
func (s *EventService) SuggestEventTimes(durationMinutes int, preferredDate time.Time) []time.Time {
	var suggestions []time.Time
	
	// Business hours: 9 AM to 5 PM
	businessStart := 9
	businessEnd := 17
	
	duration := time.Duration(durationMinutes) * time.Minute
	
	// Generate suggestions every 30 minutes during business hours
	for hour := businessStart; hour < businessEnd; hour++ {
		for minute := 0; minute < 60; minute += 30 {
			suggested := time.Date(
				preferredDate.Year(),
				preferredDate.Month(),
				preferredDate.Day(),
				hour, minute, 0, 0,
				preferredDate.Location(),
			)
			
			// Check if event would end within business hours
			endTime := suggested.Add(duration)
			if endTime.Hour() <= businessEnd {
				suggestions = append(suggestions, suggested)
			}
			
			// Limit to 10 suggestions
			if len(suggestions) >= 10 {
				break
			}
		}
		if len(suggestions) >= 10 {
			break
		}
	}
	
	return suggestions
}