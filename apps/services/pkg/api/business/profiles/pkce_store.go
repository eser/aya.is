package profiles

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrPKCEVerifierNotFound = errors.New("PKCE verifier not found or expired")

const (
	pkceVerifierBytes         = 32
	pkceVerifierExpiryMinutes = 10
)

// pkceEntry holds a code verifier with an expiry time.
type pkceEntry struct {
	ExpiresAt    time.Time
	CodeVerifier string
}

// PKCEStore stores PKCE code verifiers between OAuth initiate and callback.
type PKCEStore struct {
	verifiers map[string]*pkceEntry
	mu        sync.RWMutex
}

// NewPKCEStore creates a new PKCE store with automatic cleanup.
func NewPKCEStore() *PKCEStore {
	store := &PKCEStore{ //nolint:exhaustruct // mu zero value is valid
		verifiers: make(map[string]*pkceEntry),
	}

	go store.cleanupExpired()

	return store
}

// GeneratePKCE generates a code_verifier, stores it keyed by stateKey,
// and returns the corresponding code_challenge (S256).
func (s *PKCEStore) GeneratePKCE(stateKey string) (string, error) {
	verifierBuf := make([]byte, pkceVerifierBytes)

	_, err := rand.Read(verifierBuf)
	if err != nil {
		return "", fmt.Errorf("failed to generate PKCE verifier: %w", err)
	}

	codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBuf)

	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	s.mu.Lock()
	defer s.mu.Unlock()

	s.verifiers[stateKey] = &pkceEntry{
		CodeVerifier: codeVerifier,
		ExpiresAt:    time.Now().Add(pkceVerifierExpiryMinutes * time.Minute),
	}

	return codeChallenge, nil
}

// GetAndDelete retrieves and removes the code verifier for a given state key.
func (s *PKCEStore) GetAndDelete(stateKey string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.verifiers[stateKey]
	if !ok {
		return "", ErrPKCEVerifierNotFound
	}

	delete(s.verifiers, stateKey)

	if time.Now().After(entry.ExpiresAt) {
		return "", ErrPKCEVerifierNotFound
	}

	return entry.CodeVerifier, nil
}

// cleanupExpired periodically removes expired verifiers.
func (s *PKCEStore) cleanupExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		now := time.Now()

		for key, entry := range s.verifiers {
			if now.After(entry.ExpiresAt) {
				delete(s.verifiers, key)
			}
		}

		s.mu.Unlock()
	}
}
