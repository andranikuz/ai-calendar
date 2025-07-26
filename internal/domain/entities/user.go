package entities

import (
	"time"
)

type UserID string

type User struct {
	ID        UserID    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Profile   UserProfile `json:"profile"`
	Settings  UserSettings `json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Timezone  string `json:"timezone"`
}

type UserSettings struct {
	Language         string `json:"language"`
	DateFormat       string `json:"date_format"`
	TimeFormat       string `json:"time_format"`
	WeekStartDay     int    `json:"week_start_day"`
	NotificationEnabled bool `json:"notification_enabled"`
}