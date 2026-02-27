package resend

import "strings"

// Config holds configuration for the Resend email API.
type Config struct {
	APIKey            string `conf:"api_key"`
	FromAddress       string `conf:"from_address"       default:"bulletin@aya.is"`
	Enabled           bool   `conf:"enabled"            default:"false"`
	SandboxMode       bool   `conf:"sandbox_mode"       default:"false"`
	SandboxRecipients string `conf:"sandbox_recipients"`
}

// IsConfigured returns true if the Resend API key is set and enabled.
func (c *Config) IsConfigured() bool {
	return c.Enabled && c.APIKey != ""
}

// SandboxAllowedEmails parses the comma-separated sandbox_recipients into a set.
func (c *Config) SandboxAllowedEmails() map[string]bool {
	if c.SandboxRecipients == "" {
		return nil
	}

	parts := strings.Split(c.SandboxRecipients, ",")
	allowed := make(map[string]bool, len(parts))

	for _, email := range parts {
		trimmed := strings.TrimSpace(email)
		if trimmed != "" {
			allowed[strings.ToLower(trimmed)] = true
		}
	}

	return allowed
}
