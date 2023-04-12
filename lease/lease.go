package main

import (
	"fmt"
	"sync"
	"time"
)

// Lease represents a lease with an associated time-to-live (TTL) and expiry time.
type Lease struct {
	name         string
	ttl          time.Duration
	expiresAt    time.Time
	attachedKeys []string
}

// NewLease creates a new lease with the given name and TTL.
func NewLease(name string, ttl time.Duration) *Lease {
	return &Lease{
		name:      name,
		ttl:       ttl,
		expiresAt: time.Now().Add(ttl),
	}
}

// Refresh refreshes the expiry time of the lease to be TTL duration from the current time.
func (l *Lease) Refresh() {
	l.expiresAt = time.Now().Add(l.ttl)
}

// AttachKey attaches a key to the lease.
func (l *Lease) AttachKey(key string) {
	l.attachedKeys = append(l.attachedKeys, key)
}

// KVStore represents a key-value store with leases.
type KVStore struct {
	mu     sync.Mutex
	leases map[string]*Lease
}

// NewKVStore creates a new key-value store.
func NewKVStore() *KVStore {
	return &KVStore{
		leases: make(map[string]*Lease),
	}
}

// RegisterLease registers a new lease with the given name and TTL.
func (s *KVStore) RegisterLease(name string, ttl time.Duration) (*Lease, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the lease already exists
	if _, ok := s.leases[name]; ok {
		return nil, fmt.Errorf("lease %q already exists", name)
	}

	// Create a new lease
	lease := NewLease(name, ttl)
	s.leases[name] = lease
	return lease, nil
}

// GetLease returns the lease with the given name, or nil if it does not exist.
func (s *KVStore) GetLease(name string) *Lease {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.leases[name]
}

// RemoveLease removes the lease with the given name.
func (s *KVStore) RemoveLease(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.leases, name)
}

// RefreshLease refreshes the expiry time of the lease with the given name.
func (s *KVStore) RefreshLease(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	lease, ok := s.leases[name]
	if !ok {
		return fmt.Errorf("lease %q does not exist", name)
	}

	lease.Refresh()
	return nil
}
