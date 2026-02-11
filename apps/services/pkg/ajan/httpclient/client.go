package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Client is a drop-in replacement for http.Client with built-in circuit breaker and retry mechanisms.
type Client struct {
	*http.Client

	Config          *Config
	Transport       *ResilientTransport
	TLSClientConfig *tls.Config
	innerTransport  http.RoundTripper
	timeout         time.Duration
}

// NewClient creates a new http client with the specified circuit breaker and retry strategy.
func NewClient(options ...NewClientOption) *Client {
	client := &Client{
		Client:          nil,
		TLSClientConfig: nil,
		innerTransport:  nil,

		Config: &Config{
			CircuitBreaker: CircuitBreakerConfig{
				Enabled:               true,
				FailureThreshold:      DefaultFailureThreshold,
				ResetTimeout:          DefaultResetTimeout,
				HalfOpenSuccessNeeded: DefaultHalfOpenSuccess,
			},
			RetryStrategy: RetryStrategyConfig{
				Enabled:         true,
				MaxAttempts:     DefaultMaxAttempts,
				InitialInterval: DefaultInitialInterval,
				MaxInterval:     DefaultMaxInterval,
				Multiplier:      DefaultMultiplier,
				RandomFactor:    DefaultRandomFactor,
			},

			ServerErrorThreshold: DefaultServerErrorThreshold,
		},
		Transport: nil,
	}

	for _, option := range options {
		option(client)
	}

	if client.Transport == nil {
		var transport http.RoundTripper

		if client.innerTransport != nil {
			transport = client.innerTransport
		} else {
			// Create a copy of the default transport to avoid race conditions
			defaultTransport, transportOk := http.DefaultTransport.(*http.Transport)
			if !transportOk {
				return nil
			}

			// Clone the transport to avoid modifying the shared default transport
			clonedTransport := defaultTransport.Clone()

			// Set TLS config if provided
			if client.TLSClientConfig != nil {
				clonedTransport.TLSClientConfig = client.TLSClientConfig
			}

			transport = clonedTransport
		}

		resilientTransport := NewResilientTransport(
			transport,
			client.Config,
		)

		client.Transport = resilientTransport
	}

	client.Client = &http.Client{ //nolint:exhaustruct
		Transport: client.Transport,
		Timeout:   client.timeout,
	}

	return client
}
