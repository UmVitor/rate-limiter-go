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
	cfg := config.LoadConfig()

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
		
		memStore.StartCleanupTask(1 * time.Minute)
		store = memStore
	default:
		log.Fatalf("Unknown storage type: %s", cfg.StorageType)
	}
	defer store.Close()

	rateLimiter := limiter.NewRateLimiter(cfg, store)
	defer rateLimiter.Close()

	var limiterInterface interfaces.RateLimiter = rateLimiter
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(limiterInterface)

	router := mux.NewRouter()

	router.Use(rateLimiterMiddleware.Middleware)

	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/api/test", testHandler).Methods("GET")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Welcome to the Rate Limiter API", "status": "ok"}`)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "This is a test endpoint", "status": "ok"}`)
}
