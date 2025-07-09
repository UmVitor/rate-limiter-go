package limiter

import (
  "context"
  "testing"
  "time"

  "rate-limiter/config"
)

// MockStorage is a mock implementation of the Storage interface for testing
type MockStorage struct {
  counters      map[string]int
  blockedKeys   map[string]bool
  lastExpiration time.Duration
}

// NewMockStorage creates a new mock storage
func NewMockStorage() *MockStorage {
  return &MockStorage{
    counters:    make(map[string]int),
    blockedKeys: make(map[string]bool),
  }
}

// Get returns the current count for a key
func (m *MockStorage) Get(ctx context.Context, key string) (int, error) {
  return m.counters[key], nil
}

// Increment increments the counter for a key and returns the new value
func (m *MockStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
  m.lastExpiration = expiration
  m.counters[key]++
  return m.counters[key], nil
}

// IsBlocked checks if a key is blocked
func (m *MockStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
  return m.blockedKeys[key], nil
}

// Block blocks a key for the specified duration
func (m *MockStorage) Block(ctx context.Context, key string, duration time.Duration) error {
  m.blockedKeys[key] = true
  return nil
}

// Close closes the storage connection
func (m *MockStorage) Close() error {
  return nil
}

// TestRateLimiterIP tests the IP-based rate limiting
func TestRateLimiterIP(t *testing.T) {
  // Create a mock storage
  mockStorage := NewMockStorage()

  // Create a config with a limit of 3 requests
  cfg := &config.Config{
    IPLimit:        3,
    IPExpiration:   300,
    BlockDuration:  300,
  }

  // Create a rate limiter with the mock storage
  limiter := NewRateLimiter(cfg, mockStorage)

  // Test IP
  ip := "192.168.1.1"
  ctx := context.Background()

  // First 3 requests should be allowed
  for i := 0; i < 3; i++ {
    allowed, err := limiter.CheckIP(ctx, ip)
    if err != nil {
      t.Errorf("Error checking IP: %v", err)
    }
    if !allowed {
      t.Errorf("Request %d should be allowed", i+1)
    }
  }

  // 4th request should be blocked
  allowed, err := limiter.CheckIP(ctx, ip)
  if err != nil {
    t.Errorf("Error checking IP: %v", err)
  }
  if allowed {
    t.Error("4th request should be blocked")
  }

  // Verify the IP is now blocked
  blocked, err := mockStorage.IsBlocked(ctx, ip)
  if err != nil {
    t.Errorf("Error checking if IP is blocked: %v", err)
  }
  if !blocked {
    t.Error("IP should be blocked")
  }
}

// TestRateLimiterToken tests the token-based rate limiting
func TestRateLimiterToken(t *testing.T) {
  // Create a mock storage
  mockStorage := NewMockStorage()

  // Create a config with a limit of 5 requests for tokens
  cfg := &config.Config{
    TokenLimit:      5,
    TokenExpiration: 300,
    BlockDuration:   300,
  }

  // Create a rate limiter with the mock storage
  limiter := NewRateLimiter(cfg, mockStorage)

  // Test token
  token := "test-token"
  ctx := context.Background()

  // First 5 requests should be allowed
  for i := 0; i < 5; i++ {
    allowed, err := limiter.CheckToken(ctx, token)
    if err != nil {
      t.Errorf("Error checking token: %v", err)
    }
    if !allowed {
      t.Errorf("Request %d should be allowed", i+1)
    }
  }

  // 6th request should be blocked
  allowed, err := limiter.CheckToken(ctx, token)
  if err != nil {
    t.Errorf("Error checking token: %v", err)
  }
  if allowed {
    t.Error("6th request should be blocked")
  }

  // Verify the token is now blocked
  blocked, err := mockStorage.IsBlocked(ctx, token)
  if err != nil {
    t.Errorf("Error checking if token is blocked: %v", err)
  }
  if !blocked {
    t.Error("Token should be blocked")
  }
}
