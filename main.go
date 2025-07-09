package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"rate-limiter/config"
	"rate-limiter/interfaces"
	"rate-limiter/limiter"
	"rate-limiter/middleware"
	"rate-limiter/storage"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize storage based on configuration
	var store storage.Storage
	var err error

	switch cfg.StorageType {
	case config.StorageTypeRedis:
		log.Println("Using Redis storage")
		store, err = storage.NewRedisStorage(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize Redis storage: %v", err)
		}
	case config.StorageTypeMemory:
		log.Println("Using in-memory storage")
		memStore := storage.NewMemoryStorage()
		// Start cleanup task to remove expired items every minute
		memStore.StartCleanupTask(1 * time.Minute)
		store = memStore
	default:
		log.Fatalf("Unknown storage type: %s", cfg.StorageType)
	}
	defer store.Close()

	// Initialize rate limiter
	rateLimiter := limiter.NewRateLimiter(cfg, store)
	defer rateLimiter.Close()

	// Initialize middleware
	// Usamos a interface interfaces.RateLimiter para garantir que o middleware possa usar qualquer implementação
	var limiterInterface interfaces.RateLimiter = rateLimiter
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(limiterInterface)

	// Initialize router
	router := mux.NewRouter()

	// Apply middleware to all routes
	router.Use(rateLimiterMiddleware.Middleware)

	// Define routes
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/api/test", testHandler).Methods("GET")

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// homeHandler handles the root endpoint
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Welcome to the Rate Limiter API", "status": "ok"}`)
}

// testHandler handles the test endpoint
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "This is a test endpoint", "status": "ok"}`)
}