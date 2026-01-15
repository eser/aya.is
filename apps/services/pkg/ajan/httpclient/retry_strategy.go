package httpclient

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

const (
	DefaultMaxAttempts     = 3
	DefaultInitialInterval = 100 * time.Millisecond
	DefaultMaxInterval     = 10 * time.Second
	DefaultMultiplier      = 2.0
	DefaultRandomFactor    = 0.1
	randomNumberRange      = 1000 // Range for random number generation in jitter calculation
)

type RetryStrategy struct {
	Config *RetryStrategyConfig
}

// NewRetryStrategy creates a new retry strategy with the specified parameters.
func NewRetryStrategy(config *RetryStrategyConfig) *RetryStrategy {
	return &RetryStrategy{
		Config: config,
	}
}

func (r *RetryStrategy) NextBackoff(attempt uint) time.Duration {
	if attempt >= r.Config.MaxAttempts {
		return 0
	}

	// Calculate exponential backoff
	backoff := float64(r.Config.InitialInterval) * math.Pow(r.Config.Multiplier, float64(attempt))

	// Apply random factor
	if r.Config.RandomFactor > 0 {
		// Use crypto/rand for secure random number generation
		n, err := rand.Int(rand.Reader, big.NewInt(randomNumberRange)) //nolint:varnamelen
		if err != nil {
			// Fallback to no jitter if random generation fails
			return time.Duration(backoff)
		}

		random := 1 + r.Config.RandomFactor*(2*float64(n.Int64())/float64(randomNumberRange)-1)
		backoff *= random
	}

	// Ensure we don't exceed max interval
	if backoff > float64(r.Config.MaxInterval) {
		backoff = float64(r.Config.MaxInterval)
	}

	return time.Duration(backoff)
}
