package valueobjects

import (
	"encoding/json"
	"fmt"
)

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

func (p Priority) String() string {
	return string(p)
}

func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	default:
		return false
	}
}

// MarshalJSON implements json.Marshaler
func (p Priority) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(p))
}

// UnmarshalJSON implements json.Unmarshaler
func (p *Priority) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	
	priority := Priority(s)
	if !priority.IsValid() {
		return fmt.Errorf("invalid priority: %s", s)
	}
	
	*p = priority
	return nil
}