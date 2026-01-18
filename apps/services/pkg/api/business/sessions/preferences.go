package sessions

import (
	"errors"
	"time"
)

// Standard preference keys.
const (
	PrefKeyTheme    = "theme"
	PrefKeyLocale   = "locale"
	PrefKeyTimezone = "timezone"
)

// Valid theme values.
const (
	ThemeLight  = "light"
	ThemeDark   = "dark"
	ThemeSystem = "system"
)

var (
	ErrInvalidPreferenceKey   = errors.New("invalid preference key")
	ErrInvalidPreferenceValue = errors.New("invalid preference value")
	ErrPreferenceNotFound     = errors.New("preference not found")
	ErrFailedToGetPreference  = errors.New("failed to get preference")
	ErrFailedToSetPreference  = errors.New("failed to set preference")
)

// AllowedPreferenceKeys is the set of valid preference keys.
var AllowedPreferenceKeys = map[string]bool{
	PrefKeyTheme:    true,
	PrefKeyLocale:   true,
	PrefKeyTimezone: true,
}

// SessionPreference represents a single session preference.
type SessionPreference struct {
	SessionID string    `json:"session_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SessionPreferences is a map of preference key to value.
type SessionPreferences map[string]string

// ValidatePreferenceKey checks if a preference key is valid.
func ValidatePreferenceKey(key string) error {
	if !AllowedPreferenceKeys[key] {
		return ErrInvalidPreferenceKey
	}

	return nil
}

// ValidateThemeValue checks if a theme value is valid.
func ValidateThemeValue(value string) error {
	switch value {
	case ThemeLight, ThemeDark, ThemeSystem:
		return nil
	default:
		return ErrInvalidPreferenceValue
	}
}

// ValidatePreference validates a preference key-value pair.
func ValidatePreference(key, value string) error {
	err := ValidatePreferenceKey(key)
	if err != nil {
		return err
	}

	// Specific validation based on key
	switch key {
	case PrefKeyTheme:
		return ValidateThemeValue(value)
	case PrefKeyLocale:
		// Accept any non-empty locale for now
		if value == "" {
			return ErrInvalidPreferenceValue
		}

		return nil
	case PrefKeyTimezone:
		// Accept any non-empty timezone for now
		if value == "" {
			return ErrInvalidPreferenceValue
		}

		return nil
	default:
		return nil
	}
}
