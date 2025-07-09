package storage

import (
  "context"
  "fmt"
  "time"

  "github.com/go-redis/redis/v8"
  "rate-limiter/config"
)

// RedisStorage implements the Storage interface using Redis
type RedisStorage struct {
  client *redis.Client
}

// NewRedisStorage creates a new Redis storage instance
func NewRedisStorage(cfg *config.Config) (*RedisStorage, error) {
  client := redis.NewClient(&redis.Options{
    Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
    Password: cfg.RedisPassword,
    DB:       cfg.RedisDB,
  })

  // Test the connection
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  if _, err := client.Ping(ctx).Result(); err != nil {
    return nil, fmt.Errorf("failed to connect to Redis: %w", err)
  }

  return &RedisStorage{
    client: client,
  }, nil
}

// Get returns the current count for a key
func (s *RedisStorage) Get(ctx context.Context, key string) (int, error) {
  val, err := s.client.Get(ctx, key).Int()
  if err == redis.Nil {
    return 0, nil
  }
  return val, err
}

// Increment increments the counter for a key and returns the new value
func (s *RedisStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
  pipe := s.client.Pipeline()
  incr := pipe.Incr(ctx, key)
  pipe.Expire(ctx, key, expiration)
  _, err := pipe.Exec(ctx)
  if err != nil {
    return 0, err
  }
  return int(incr.Val()), nil
}

// IsBlocked checks if a key is blocked
func (s *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
  blockedKey := fmt.Sprintf("blocked:%s", key)
  exists, err := s.client.Exists(ctx, blockedKey).Result()
  if err != nil {
    return false, err
  }
  return exists > 0, nil
}

// Block blocks a key for the specified duration
func (s *RedisStorage) Block(ctx context.Context, key string, duration time.Duration) error {
  blockedKey := fmt.Sprintf("blocked:%s", key)
  return s.client.Set(ctx, blockedKey, 1, duration).Err()
}

// Close closes the Redis connection
func (s *RedisStorage) Close() error {
  return s.client.Close()
}
