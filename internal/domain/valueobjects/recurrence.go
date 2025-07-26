package valueobjects

import (
	"fmt"
	"strings"
	"time"
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

// NewRecurrenceRule creates a new recurrence rule with validation
func NewRecurrenceRule(frequency Frequency, interval int) *RecurrenceRule {
	if interval <= 0 {
		interval = 1
	}
	
	return &RecurrenceRule{
		Frequency: frequency,
		Interval:  interval,
	}
}

// SetUntil sets the end date for the recurrence
func (r *RecurrenceRule) SetUntil(until time.Time) *RecurrenceRule {
	r.Until = &until
	r.Count = nil // Clear count if until is set
	return r
}

// SetCount sets the number of occurrences
func (r *RecurrenceRule) SetCount(count int) *RecurrenceRule {
	r.Count = &count
	r.Until = nil // Clear until if count is set
	return r
}

// SetByDay sets the weekdays for weekly recurrence
func (r *RecurrenceRule) SetByDay(weekdays ...Weekday) *RecurrenceRule {
	r.ByDay = weekdays
	return r
}

// SetByMonth sets the months for yearly recurrence
func (r *RecurrenceRule) SetByMonth(months ...Month) *RecurrenceRule {
	r.ByMonth = months
	return r
}

// IsValid validates the recurrence rule
func (r *RecurrenceRule) IsValid() bool {
	if r.Interval <= 0 {
		return false
	}
	
	validFrequencies := []Frequency{FrequencyDaily, FrequencyWeekly, FrequencyMonthly, FrequencyYearly}
	validFreq := false
	for _, f := range validFrequencies {
		if r.Frequency == f {
			validFreq = true
			break
		}
	}
	
	if !validFreq {
		return false
	}
	
	// Cannot have both until and count
	if r.Until != nil && r.Count != nil {
		return false
	}
	
	// Validate weekdays for weekly recurrence
	if r.Frequency == FrequencyWeekly && len(r.ByDay) > 0 {
		validWeekdays := []Weekday{WeekdayMonday, WeekdayTuesday, WeekdayWednesday, WeekdayThursday, WeekdayFriday, WeekdaySaturday, WeekdaySunday}
		for _, wd := range r.ByDay {
			valid := false
			for _, vwd := range validWeekdays {
				if wd == vwd {
					valid = true
					break
				}
			}
			if !valid {
				return false
			}
		}
	}
	
	// Validate months for yearly recurrence
	if r.Frequency == FrequencyYearly && len(r.ByMonth) > 0 {
		for _, m := range r.ByMonth {
			if m < January || m > December {
				return false
			}
		}
	}
	
	return true
}

// ToRRULE converts the recurrence rule to RFC 5545 RRULE format
func (r *RecurrenceRule) ToRRULE() string {
	if !r.IsValid() {
		return ""
	}
	
	parts := []string{
		fmt.Sprintf("FREQ=%s", r.Frequency),
	}
	
	if r.Interval > 1 {
		parts = append(parts, fmt.Sprintf("INTERVAL=%d", r.Interval))
	}
	
	if r.Until != nil {
		parts = append(parts, fmt.Sprintf("UNTIL=%s", r.Until.UTC().Format("20060102T150405Z")))
	}
	
	if r.Count != nil {
		parts = append(parts, fmt.Sprintf("COUNT=%d", *r.Count))
	}
	
	if len(r.ByDay) > 0 {
		weekdays := make([]string, len(r.ByDay))
		for i, wd := range r.ByDay {
			weekdays[i] = string(wd)
		}
		parts = append(parts, fmt.Sprintf("BYDAY=%s", strings.Join(weekdays, ",")))
	}
	
	if len(r.ByMonth) > 0 {
		months := make([]string, len(r.ByMonth))
		for i, m := range r.ByMonth {
			months[i] = fmt.Sprintf("%d", m)
		}
		parts = append(parts, fmt.Sprintf("BYMONTH=%s", strings.Join(months, ",")))
	}
	
	return strings.Join(parts, ";")
}

// GetNextOccurrence calculates the next occurrence after the given time
func (r *RecurrenceRule) GetNextOccurrence(after time.Time) time.Time {
	if !r.IsValid() {
		return time.Time{}
	}
	
	next := after
	
	switch r.Frequency {
	case FrequencyDaily:
		next = next.AddDate(0, 0, r.Interval)
	case FrequencyWeekly:
		next = next.AddDate(0, 0, 7*r.Interval)
	case FrequencyMonthly:
		next = next.AddDate(0, r.Interval, 0)
	case FrequencyYearly:
		next = next.AddDate(r.Interval, 0, 0)
	}
	
	return next
}

// Common recurrence patterns
func DailyRecurrence() *RecurrenceRule {
	return NewRecurrenceRule(FrequencyDaily, 1)
}

func WeeklyRecurrence() *RecurrenceRule {
	return NewRecurrenceRule(FrequencyWeekly, 1)
}

func WeekdaysRecurrence() *RecurrenceRule {
	return NewRecurrenceRule(FrequencyWeekly, 1).SetByDay(
		WeekdayMonday, WeekdayTuesday, WeekdayWednesday, WeekdayThursday, WeekdayFriday,
	)
}

func MonthlyRecurrence() *RecurrenceRule {
	return NewRecurrenceRule(FrequencyMonthly, 1)
}

func YearlyRecurrence() *RecurrenceRule {
	return NewRecurrenceRule(FrequencyYearly, 1)
}