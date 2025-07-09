package limiter

import (
  "context"
  "fmt"
  "time"

  "rate-limiter/config"
  "rate-limiter/interfaces"
  "rate-limiter/storage"
)

// Ensure RateLimiter implements the interfaces.RateLimiter interface
var _ interfaces.RateLimiter = (*RateLimiter)(nil)

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
  storage       storage.Storage
  ipLimit       int
  ipExpiration  time.Duration
  tokenLimit    int
  tokenExpiration time.Duration
  blockDuration time.Duration
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(cfg *config.Config, store storage.Storage) *RateLimiter {
  return &RateLimiter{
    storage:         store,
    ipLimit:         cfg.IPLimit,
    ipExpiration:    time.Duration(cfg.IPExpiration) * time.Second,
    tokenLimit:      cfg.TokenLimit,
    tokenExpiration: time.Duration(cfg.TokenExpiration) * time.Second,
    blockDuration:   time.Duration(cfg.BlockDuration) * time.Second,
  }
}

// CheckIP checks if an IP address has exceeded its rate limit
func (rl *RateLimiter) CheckIP(ctx context.Context, ip string) (bool, error) {
  // Check if the IP is blocked
  blocked, err := rl.storage.IsBlocked(ctx, ip)
  if err != nil {
    return false, err
  }
  if blocked {
    return false, nil
  }

  // Get the current count for this IP
  key := fmt.Sprintf("ip:%s", ip)
  count, err := rl.storage.Increment(ctx, key, rl.ipExpiration)
  if err != nil {
    return false, err
  }

  // If the count exceeds the limit, block the IP
  if count > rl.ipLimit {
    if err := rl.storage.Block(ctx, ip, rl.blockDuration); err != nil {
      return false, err
    }
    return false, nil
  }

  return true, nil
}

// CheckToken checks if a token has exceeded its rate limit
func (rl *RateLimiter) CheckToken(ctx context.Context, token string) (bool, error) {
  // Check if the token is blocked
  blocked, err := rl.storage.IsBlocked(ctx, token)
  if err != nil {
    return false, err
  }
  if blocked {
    return false, nil
  }

  // Get the current count for this token
  key := fmt.Sprintf("token:%s", token)
  count, err := rl.storage.Increment(ctx, key, rl.tokenExpiration)
  if err != nil {
    return false, err
  }

  // If the count exceeds the limit, block the token
  if count > rl.tokenLimit {
    if err := rl.storage.Block(ctx, token, rl.blockDuration); err != nil {
      return false, err
    }
    return false, nil
  }

  return true, nil
}

// Close closes the rate limiter and its storage
func (rl *RateLimiter) Close() error {
  return rl.storage.Close()
}
