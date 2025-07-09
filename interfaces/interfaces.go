package interfaces

import (
  "context"
)

// RateLimiter defines the interface for rate limiters
type RateLimiter interface {
  // CheckIP checks if an IP address has exceeded its rate limit
  CheckIP(ctx context.Context, ip string) (bool, error)

  // CheckToken checks if a token has exceeded its rate limit
  CheckToken(ctx context.Context, token string) (bool, error)

  // Close closes the rate limiter
  Close() error
}
