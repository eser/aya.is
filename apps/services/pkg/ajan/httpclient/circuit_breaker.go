package httpclient

import (
	"sync"
	"time"
)

//go:generate go tool stringer -type CircuitState -trimprefix CircuitState
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateHalfOpen
	StateOpen
)

type CircuitBreaker struct {
	lastFailureTime time.Time

	Config *CircuitBreakerConfig

	state                CircuitState
	failureCount         uint
	halfOpenSuccessCount uint
	mu                   sync.RWMutex
}

func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{ //nolint:exhaustruct
		Config: config,
		state:  StateClosed,
	}
}

func (cb *CircuitBreaker) IsAllowed() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.Config.ResetTimeout {
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.state = StateHalfOpen
			cb.halfOpenSuccessCount = 0
			cb.mu.Unlock()
			cb.mu.RLock()

			return true
		}

		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) OnSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		cb.halfOpenSuccessCount++
		if cb.halfOpenSuccessCount >= cb.Config.HalfOpenSuccessNeeded {
			cb.state = StateClosed
			cb.failureCount = 0
		}
	}

	if cb.state == StateClosed {
		cb.failureCount = 0
	}
}

func (cb *CircuitBreaker) OnFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.state == StateHalfOpen ||
		(cb.state == StateClosed && cb.failureCount >= cb.Config.FailureThreshold) {
		cb.state = StateOpen
	}
}

func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return cb.state
}
