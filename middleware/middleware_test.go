package middleware

import (
  "context"
  "net/http"
  "net/http/httptest"
  "testing"

  "rate-limiter/interfaces"
)

// MockRateLimiter é uma implementação mock da interface interfaces.RateLimiter para testes
type MockRateLimiter struct {
  allowIP    bool
  allowToken bool
  err        error
}

// Garantir que MockRateLimiter implementa a interface interfaces.RateLimiter
var _ interfaces.RateLimiter = (*MockRateLimiter)(nil)

// CheckIP mocks the IP check
func (m *MockRateLimiter) CheckIP(ctx context.Context, ip string) (bool, error) {
  return m.allowIP, m.err
}

// CheckToken mocks the token check
func (m *MockRateLimiter) CheckToken(ctx context.Context, token string) (bool, error) {
  return m.allowToken, m.err
}

// Close mocks the close method
func (m *MockRateLimiter) Close() error {
  return nil
}

// TestMiddlewareIPAllowed tests that requests with allowed IPs pass through
func TestMiddlewareIPAllowed(t *testing.T) {
  // Create a mock rate limiter that allows all IPs
  mockLimiter := &MockRateLimiter{
    allowIP: true,
  }

  // Create the middleware
  middleware := NewRateLimiterMiddleware(mockLimiter)

  // Create a test handler that will be called if the middleware allows the request
  testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
  })

  // Create a test request
  req := httptest.NewRequest("GET", "/test", nil)
  req.RemoteAddr = "192.168.1.1:12345"

  // Create a response recorder
  rr := httptest.NewRecorder()

  // Apply the middleware to the test handler and serve the request
  middleware.Middleware(testHandler).ServeHTTP(rr, req)

  // Check the response status code
  if status := rr.Code; status != http.StatusOK {
    t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
  }
}

// TestMiddlewareIPBlocked tests that requests with blocked IPs are rejected
func TestMiddlewareIPBlocked(t *testing.T) {
  // Create a mock rate limiter that blocks all IPs
  mockLimiter := &MockRateLimiter{
    allowIP: false,
  }

  // Create the middleware
  middleware := NewRateLimiterMiddleware(mockLimiter)

  // Create a test handler that should not be called
  testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    t.Error("Handler should not be called for blocked IP")
  })

  // Create a test request
  req := httptest.NewRequest("GET", "/test", nil)
  req.RemoteAddr = "192.168.1.1:12345"

  // Create a response recorder
  rr := httptest.NewRecorder()

  // Apply the middleware to the test handler and serve the request
  middleware.Middleware(testHandler).ServeHTTP(rr, req)

  // Check the response status code
  if status := rr.Code; status != http.StatusTooManyRequests {
    t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
  }
}

// TestMiddlewareTokenAllowed tests that requests with allowed tokens pass through
func TestMiddlewareTokenAllowed(t *testing.T) {
  // Create a mock rate limiter that allows all tokens
  mockLimiter := &MockRateLimiter{
    allowToken: true,
  }

  // Create the middleware
  middleware := NewRateLimiterMiddleware(mockLimiter)

  // Create a test handler that will be called if the middleware allows the request
  testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
  })

  // Create a test request with a token
  req := httptest.NewRequest("GET", "/test", nil)
  req.Header.Set(TokenHeader, "test-token")

  // Create a response recorder
  rr := httptest.NewRecorder()

  // Apply the middleware to the test handler and serve the request
  middleware.Middleware(testHandler).ServeHTTP(rr, req)

  // Check the response status code
  if status := rr.Code; status != http.StatusOK {
    t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
  }
}

// TestMiddlewareTokenBlocked tests that requests with blocked tokens are rejected
func TestMiddlewareTokenBlocked(t *testing.T) {
  // Create a mock rate limiter that blocks all tokens
  mockLimiter := &MockRateLimiter{
    allowToken: false,
  }

  // Create the middleware
  middleware := NewRateLimiterMiddleware(mockLimiter)

  // Create a test handler that should not be called
  testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    t.Error("Handler should not be called for blocked token")
  })

  // Create a test request with a token
  req := httptest.NewRequest("GET", "/test", nil)
  req.Header.Set(TokenHeader, "test-token")

  // Create a response recorder
  rr := httptest.NewRecorder()

  // Apply the middleware to the test handler and serve the request
  middleware.Middleware(testHandler).ServeHTTP(rr, req)

  // Check the response status code
  if status := rr.Code; status != http.StatusTooManyRequests {
    t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
  }
}
