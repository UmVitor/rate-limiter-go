package middleware

import (
  "encoding/json"
  "net"
  "net/http"
  "strings"

  "rate-limiter/interfaces"
)

const (
  // TokenHeader is the header name for API token
  TokenHeader = "API_KEY"
)

// RateLimiterMiddleware is a middleware that limits request rates
type RateLimiterMiddleware struct {
  limiter interfaces.RateLimiter
}

// NewRateLimiterMiddleware creates a new rate limiter middleware
func NewRateLimiterMiddleware(limiter interfaces.RateLimiter) *RateLimiterMiddleware {
  return &RateLimiterMiddleware{
    limiter: limiter,
  }
}

// Middleware returns a handler function that implements rate limiting
func (m *RateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Check if a token is provided
    token := r.Header.Get(TokenHeader)
    if token != "" {
      // Token-based rate limiting takes precedence
      allowed, err := m.limiter.CheckToken(ctx, token)
      if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
      }
      if !allowed {
        sendRateLimitExceededResponse(w)
        return
      }
    } else {
      // IP-based rate limiting
      ip := getClientIP(r)
      allowed, err := m.limiter.CheckIP(ctx, ip)
      if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
      }
      if !allowed {
        sendRateLimitExceededResponse(w)
        return
      }
    }

    // If we get here, the request is allowed
    next.ServeHTTP(w, r)
  })
}

// Helper function to get the client's IP address
func getClientIP(r *http.Request) string {
  // Check for X-Forwarded-For header
  xForwardedFor := r.Header.Get("X-Forwarded-For")
  if xForwardedFor != "" {
    // X-Forwarded-For can contain multiple IPs, take the first one
    ips := strings.Split(xForwardedFor, ",")
    if len(ips) > 0 {
      return strings.TrimSpace(ips[0])
    }
  }

  // If no X-Forwarded-For header, use RemoteAddr
  ip, _, err := net.SplitHostPort(r.RemoteAddr)
  if err != nil {
    // If there's an error, just return the RemoteAddr as is
    return r.RemoteAddr
  }
  return ip
}

// Helper function to send a rate limit exceeded response
func sendRateLimitExceededResponse(w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusTooManyRequests) // 429 Too Many Requests

  response := map[string]string{
    "error":   "Rate limit exceeded",
    "message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
  }

  json.NewEncoder(w).Encode(response)
}
