package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// ValidateUserCreation validates user data before creation
func (s *UserService) ValidateUserCreation(user *entities.User) error {
	if err := s.ValidateEmail(user.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	
	if err := s.ValidateName(user.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	
	if err := s.ValidateProfile(&user.Profile); err != nil {
		return fmt.Errorf("invalid profile: %w", err)
	}
	
	if err := s.ValidateSettings(&user.Settings); err != nil {
		return fmt.Errorf("invalid settings: %w", err)
	}
	
	return nil
}

// ValidateEmail validates email format
func (s *UserService) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	
	if len(email) > 255 {
		return fmt.Errorf("email is too long (max 255 characters)")
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// ValidateName validates user name
func (s *UserService) ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	
	if len(name) < 2 {
		return fmt.Errorf("name is too short (min 2 characters)")
	}
	
	if len(name) > 100 {
		return fmt.Errorf("name is too long (max 100 characters)")
	}
	
	// Check for invalid characters
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-\.]+$`)
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("name contains invalid characters")
	}
	
	return nil
}

// ValidateProfile validates user profile data
func (s *UserService) ValidateProfile(profile *entities.UserProfile) error {
	if profile.FirstName != "" && len(profile.FirstName) > 50 {
		return fmt.Errorf("first name is too long (max 50 characters)")
	}
	
	if profile.LastName != "" && len(profile.LastName) > 50 {
		return fmt.Errorf("last name is too long (max 50 characters)")
	}
	
	if profile.Timezone != "" {
		if err := s.ValidateTimezone(profile.Timezone); err != nil {
			return fmt.Errorf("invalid timezone: %w", err)
		}
	}
	
	if profile.Avatar != "" {
		if err := s.ValidateAvatarURL(profile.Avatar); err != nil {
			return fmt.Errorf("invalid avatar URL: %w", err)
		}
	}
	
	return nil
}

// ValidateSettings validates user settings
func (s *UserService) ValidateSettings(settings *entities.UserSettings) error {
	if settings.Language != "" {
		if err := s.ValidateLanguage(settings.Language); err != nil {
			return fmt.Errorf("invalid language: %w", err)
		}
	}
	
	if settings.DateFormat != "" {
		if err := s.ValidateDateFormat(settings.DateFormat); err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
	}
	
	if settings.TimeFormat != "" {
		if err := s.ValidateTimeFormat(settings.TimeFormat); err != nil {
			return fmt.Errorf("invalid time format: %w", err)
		}
	}
	
	if settings.WeekStartDay < 0 || settings.WeekStartDay > 6 {
		return fmt.Errorf("invalid week start day (must be 0-6)")
	}
	
	return nil
}

// ValidateTimezone validates timezone string
func (s *UserService) ValidateTimezone(tz string) error {
	_, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("unknown timezone: %s", tz)
	}
	return nil
}

// ValidateAvatarURL validates avatar URL format
func (s *UserService) ValidateAvatarURL(url string) error {
	if len(url) > 500 {
		return fmt.Errorf("avatar URL is too long (max 500 characters)")
	}
	
	urlRegex := regexp.MustCompile(`^https?://[^\s]+\.(jpg|jpeg|png|gif|webp)(\?[^\s]*)?$`)
	if !urlRegex.MatchString(url) {
		return fmt.Errorf("invalid avatar URL format")
	}
	
	return nil
}

// ValidateLanguage validates language code
func (s *UserService) ValidateLanguage(lang string) error {
	validLanguages := []string{"en", "es", "fr", "de", "it", "pt", "ru", "zh", "ja", "ko"}
	
	for _, validLang := range validLanguages {
		if lang == validLang {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported language: %s", lang)
}

// ValidateDateFormat validates date format string
func (s *UserService) ValidateDateFormat(format string) error {
	validFormats := []string{"YYYY-MM-DD", "MM/DD/YYYY", "DD/MM/YYYY", "DD.MM.YYYY"}
	
	for _, validFormat := range validFormats {
		if format == validFormat {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported date format: %s", format)
}

// ValidateTimeFormat validates time format string
func (s *UserService) ValidateTimeFormat(format string) error {
	validFormats := []string{"12h", "24h"}
	
	for _, validFormat := range validFormats {
		if format == validFormat {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported time format: %s", format)
}

// SanitizeName sanitizes user name by trimming spaces and normalizing
func (s *UserService) SanitizeName(name string) string {
	// Trim spaces
	name = strings.TrimSpace(name)
	
	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	name = spaceRegex.ReplaceAllString(name, " ")
	
	return name
}

// SanitizeEmail sanitizes email by trimming and lowercasing
func (s *UserService) SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// CanDeleteUser checks if user can be deleted (business rules)
func (s *UserService) CanDeleteUser(ctx context.Context, userID entities.UserID) error {
	// Add business logic here, e.g.:
	// - Check if user has active subscriptions
	// - Check if user has pending payments
	// - Check if user is an admin of any organization
	
	// For now, allow deletion
	return nil
}

// GetDefaultUserSettings returns default settings for new users
func (s *UserService) GetDefaultUserSettings() entities.UserSettings {
	return entities.UserSettings{
		Language:            "en",
		DateFormat:          "YYYY-MM-DD",
		TimeFormat:          "24h",
		WeekStartDay:        1, // Monday
		NotificationEnabled: true,
	}
}

// GetDefaultUserProfile returns default profile for new users
func (s *UserService) GetDefaultUserProfile() entities.UserProfile {
	return entities.UserProfile{
		Timezone: "UTC",
	}
}