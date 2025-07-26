package entities

import "github.com/andranikuz/smart-goal-calendar/internal/domain/valueobjects"

type Priority = valueobjects.Priority

const (
	PriorityLow      = valueobjects.PriorityLow
	PriorityMedium   = valueobjects.PriorityMedium
	PriorityHigh     = valueobjects.PriorityHigh
	PriorityCritical = valueobjects.PriorityCritical
)