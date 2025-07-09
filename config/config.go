package config

import (
  "log"
  "os"
  "strconv"

  "github.com/joho/godotenv"
)

// StorageType defines the type of storage to use
type StorageType string

const (
  // StorageTypeRedis uses Redis for storage
  StorageTypeRedis StorageType = "redis"
  // StorageTypeMemory uses in-memory storage
  StorageTypeMemory StorageType = "memory"
)

// Config holds all configuration for the application
type Config struct {
  // Rate limiter configuration
  IPLimit           int
  IPExpiration      int
  TokenLimit        int
  TokenExpiration   int
  BlockDuration     int

  // Storage configuration
  StorageType StorageType

  // Redis configuration
  RedisHost     string
  RedisPort     string
  RedisPassword string
  RedisDB       int

  // Server configuration
  ServerPort string
}

// LoadConfig loads the configuration from environment variables or .env file
func LoadConfig() *Config {
  // Load .env file if it exists
  _ = godotenv.Load()

  // Determine storage type
  storageType := StorageType(getEnv("STORAGE_TYPE", string(StorageTypeRedis)))
  if storageType != StorageTypeRedis && storageType != StorageTypeMemory {
    log.Printf("Warning: Invalid storage type '%s', using Redis", storageType)
    storageType = StorageTypeRedis
  }

  return &Config{
    // Rate limiter configuration
    IPLimit:         getEnvAsInt("RATE_LIMITER_IP_LIMIT", 10),
    IPExpiration:    getEnvAsInt("RATE_LIMITER_IP_EXPIRATION", 300),
    TokenLimit:      getEnvAsInt("RATE_LIMITER_TOKEN_LIMIT", 100),
    TokenExpiration: getEnvAsInt("RATE_LIMITER_TOKEN_EXPIRATION", 300),
    BlockDuration:   getEnvAsInt("RATE_LIMITER_BLOCK_DURATION", 300),

    // Storage configuration
    StorageType: storageType,

    // Redis configuration
    RedisHost:     getEnv("REDIS_HOST", "localhost"),
    RedisPort:     getEnv("REDIS_PORT", "6379"),
    RedisPassword: getEnv("REDIS_PASSWORD", ""),
    RedisDB:       getEnvAsInt("REDIS_DB", 0),

    // Server configuration
    ServerPort: getEnv("SERVER_PORT", "8080"),
  }
}

// Helper function to get an environment variable or return a default value
func getEnv(key, defaultValue string) string {
  if value, exists := os.LookupEnv(key); exists {
    return value
  }
  return defaultValue
}

// Helper function to get an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
  if valueStr, exists := os.LookupEnv(key); exists {
    if value, err := strconv.Atoi(valueStr); err == nil {
      return value
    } else {
      log.Printf("Warning: Invalid value for %s, using default: %d", key, defaultValue)
    }
  }
  return defaultValue
}
