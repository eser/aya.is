package sessions

// Config holds configuration for the sessions module.
type Config struct {
	RateLimit RateLimitConfig `conf:"rate_limit"`
}

// RateLimitConfig holds rate limiting configuration for session creation.
type RateLimitConfig struct {
	PerIP  int `conf:"per_ip" default:"10"`  // Max sessions per IP per hour
	Global int `conf:"global" default:"100"` // Max sessions globally per minute
}
