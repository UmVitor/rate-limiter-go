package storage

import (
	"context"
	"sync"
	"time"
)

// Item represents a stored item with expiration time
type Item struct {
	Value      int
	Expiration time.Time
}

// MemoryStorage implements the Storage interface using in-memory storage
type MemoryStorage struct {
	counters    map[string]*Item
	blockedKeys map[string]time.Time
	mutex       sync.RWMutex
}

// NewMemoryStorage creates a new memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		counters:    make(map[string]*Item),
		blockedKeys: make(map[string]time.Time),
	}
}

// Get returns the current count for a key
func (s *MemoryStorage) Get(ctx context.Context, key string) (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	item, exists := s.counters[key]
	if !exists {
		return 0, nil
	}

	// Check if the item has expired
	if time.Now().After(item.Expiration) {
		return 0, nil
	}

	return item.Value, nil
}

// Increment increments the counter for a key and returns the new value
func (s *MemoryStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if the key exists and is not expired
	item, exists := s.counters[key]
	if !exists || time.Now().After(item.Expiration) {
		// Create a new item or reset an expired one
		s.counters[key] = &Item{
			Value:      1,
			Expiration: time.Now().Add(expiration),
		}
		return 1, nil
	}

	// Increment the existing item
	item.Value++
	item.Expiration = time.Now().Add(expiration)
	return item.Value, nil
}

// IsBlocked checks if a key is blocked
func (s *MemoryStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	expirationTime, exists := s.blockedKeys[key]
	if !exists {
		return false, nil
	}

	// Check if the block has expired
	if time.Now().After(expirationTime) {
		return false, nil
	}

	return true, nil
}

// Block blocks a key for the specified duration
func (s *MemoryStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.blockedKeys[key] = time.Now().Add(duration)
	return nil
}

// Close closes the storage connection (no-op for memory storage)
func (s *MemoryStorage) Close() error {
	return nil
}

// StartCleanupTask starts a background task to clean up expired items
func (s *MemoryStorage) StartCleanupTask(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			s.cleanup()
		}
	}()
}

// cleanup removes expired items from storage
func (s *MemoryStorage) cleanup() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()

	// Clean up expired counters
	for key, item := range s.counters {
		if now.After(item.Expiration) {
			delete(s.counters, key)
		}
	}

	// Clean up expired blocks
	for key, expiration := range s.blockedKeys {
		if now.After(expiration) {
			delete(s.blockedKeys, key)
		}
	}
}
