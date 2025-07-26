package entities

import (
	"time"
)

type EventID string

type Event struct {
	ID          EventID   `json:"id"`
	UserID      UserID    `json:"user_id"`
	GoalID      *GoalID   `json:"goal_id,omitempty"` // Optional связь с целью
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Timezone    string    `json:"timezone"`
	Recurrence  *RecurrenceRule `json:"recurrence,omitempty"`
	Location    string `json:"location"`
	Attendees   []Attendee `json:"attendees"`
	Status      EventStatus `json:"status"`
	ExternalID  string    `json:"external_id,omitempty"` // For Google Calendar sync
	ExternalSource string `json:"external_source,omitempty"` // 'google', 'outlook', etc.
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventStatus string

const (
	EventStatusTentative EventStatus = "tentative"
	EventStatusConfirmed EventStatus = "confirmed"
	EventStatusCancelled EventStatus = "cancelled"
)

type Location struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type Attendee struct {
	Email        string          `json:"email"`
	Name         string          `json:"name"`
	Status       AttendeeStatus  `json:"status"`
	ResponseTime *time.Time      `json:"response_time,omitempty"`
}

type AttendeeStatus string

const (
	AttendeeStatusPending   AttendeeStatus = "pending"
	AttendeeStatusAccepted  AttendeeStatus = "accepted"
	AttendeeStatusDeclined  AttendeeStatus = "declined"
	AttendeeStatusTentative AttendeeStatus = "tentative"
)

type RecurrenceRule struct {
	Frequency Frequency  `json:"frequency"`
	Interval  int        `json:"interval"`       // Every N frequency
	Until     *time.Time `json:"until,omitempty"`
	Count     *int       `json:"count,omitempty"`
	ByDay     []Weekday  `json:"by_day,omitempty"`
	ByMonth   []Month    `json:"by_month,omitempty"`
}

type Frequency string

const (
	FrequencyDaily   Frequency = "DAILY"
	FrequencyWeekly  Frequency = "WEEKLY"
	FrequencyMonthly Frequency = "MONTHLY"
	FrequencyYearly  Frequency = "YEARLY"
)

type Weekday string

const (
	WeekdayMonday    Weekday = "MO"
	WeekdayTuesday   Weekday = "TU"
	WeekdayWednesday Weekday = "WE"
	WeekdayThursday  Weekday = "TH"
	WeekdayFriday    Weekday = "FR"
	WeekdaySaturday  Weekday = "SA"
	WeekdaySunday    Weekday = "SU"
)

type Month int

const (
	January Month = iota + 1
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

// Validation methods
func (e *Event) IsValid() bool {
	return e.Title != "" && 
		   e.UserID != "" && 
		   !e.StartTime.IsZero() && 
		   !e.EndTime.IsZero() && 
		   e.StartTime.Before(e.EndTime)
}

func (e *Event) Duration() time.Duration {
	return e.EndTime.Sub(e.StartTime)
}

func (e *Event) IsRecurring() bool {
	return e.Recurrence != nil
}

func (e *Event) IsLinkedToGoal() bool {
	return e.GoalID != nil
}