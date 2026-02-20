package profiles

import (
	"sync"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
)

// PendingConnectionStore stores temporary OAuth connections awaiting account selection.
type PendingConnectionStore struct {
	mu          sync.RWMutex
	connections map[string]*PendingOAuthConnection
}

// NewPendingConnectionStore creates a new pending connection store.
func NewPendingConnectionStore() *PendingConnectionStore {
	store := &PendingConnectionStore{
		connections: make(map[string]*PendingOAuthConnection),
	}

	// Start cleanup goroutine
	go store.cleanupExpired()

	return store
}

// Store stores a pending connection and returns its ID.
func (s *PendingConnectionStore) Store(conn *PendingOAuthConnection) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := lib.IDsGenerateUnique()
	conn.ExpiresAt = time.Now().Add(10 * time.Minute)
	s.connections[id] = conn

	return id
}

// Get retrieves a pending connection by ID.
func (s *PendingConnectionStore) Get(id string) *PendingOAuthConnection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conn, ok := s.connections[id]
	if !ok {
		return nil
	}

	if time.Now().After(conn.ExpiresAt) {
		return nil
	}

	return conn
}

// Delete removes a pending connection.
func (s *PendingConnectionStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.connections, id)
}

// cleanupExpired periodically removes expired connections.
func (s *PendingConnectionStore) cleanupExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		now := time.Now()

		for id, conn := range s.connections {
			if now.After(conn.ExpiresAt) {
				delete(s.connections, id)
			}
		}

		s.mu.Unlock()
	}
}
