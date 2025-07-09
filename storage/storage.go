package storage

import (
  "context"
  "time"
)

// Storage defines the interface for rate limiter storage implementations
type Storage interface {
  // Get returns the current count for a key
  Get(ctx context.Context, key string) (int, error)

  // Increment increments the counter for a key and returns the new value
  Increment(ctx context.Context, key string, expiration time.Duration) (int, error)

  // IsBlocked checks if a key is blocked
  IsBlocked(ctx context.Context, key string) (bool, error)

  // Block blocks a key for the specified duration
  Block(ctx context.Context, key string, duration time.Duration) error

  // Close closes the storage connection
  Close() error
}
