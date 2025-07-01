package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/internal/config"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/internal/db"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/internal/handler"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/internal/repository"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/internal/service"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/pkg/logger"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/pkg/validator"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize logger
	log := logger.New()

	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found")
	}
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Connect to database
	database, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	// Test database connection
	if err := database.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")
	}

	log.Info().Msg("Database connection established")

	// Initialize dependencies
	queries := db.New(database)
	validator := validator.New()

	// Initialize repositories
	userRepo := repository.NewUserRepository(queries)

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, validator, log)

	// Setup routes
	router := setupRoutes(userHandler)

	// Setup server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Str("address", server.Addr).Msg("Starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

func setupRoutes(userHandler *handler.UserHandler) *mux.Router {
	router := mux.NewRouter()

	// API versioning
	api := router.PathPrefix("/api/v1").Subrouter()

	// User routes
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Auth routes
	api.HandleFunc("/auth/login", userHandler.Login).Methods("POST")

	// Add CORS middleware
	router.Use(corsMiddleware)

	// Add logging middleware
	router.Use(loggingMiddleware)

	return router
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request
		duration := time.Since(start)

		log := logger.New()
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Dur("duration", duration).
			Msg("HTTP request")
	})
}
