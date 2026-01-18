package protection

import "time"

// Config holds configuration for the protection module.
type Config struct {
	POWChallenge POWChallengeConfig `conf:"pow_challenge"`
}

// POWChallengeConfig holds configuration for PoW challenges.
type POWChallengeConfig struct {
	Enabled    bool          `conf:"enabled"    default:"true"`
	Difficulty int           `conf:"difficulty" default:"16"`  // Number of leading zero bits
	Expiry     time.Duration `conf:"expiry"     default:"30s"` // Challenge validity period
}
