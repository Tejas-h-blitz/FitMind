package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fitmind/handlers"
	"fitmind/middleware"
	"fitmind/services"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load environment variables (.env) if it exists
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Services
	sbService := services.NewSupabaseService()
	qdService := services.NewQdrantService()
	ragService := services.NewRAGService()

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(sbService)
	docHandler := handlers.NewDocumentsHandler(sbService, qdService, ragService)
	chatHandler := handlers.NewChatHandler(sbService, ragService)
	healthHandler := handlers.NewHealthHandler(sbService)
	analyzerHandler := handlers.NewAnalyzerHandler(sbService)
	mealPlanHandler := handlers.NewMealPlanHandler(sbService, ragService)
	workoutPlanHandler := handlers.NewWorkoutPlanHandler(sbService, ragService)

	// Router Setup
	r := chi.NewRouter()

	// Global Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS)

	// Root Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy","service":"backend"}`))
	})

	// API Routes (Protected)
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Auth)

		// Auth
		r.Post("/auth/verify", authHandler.Verify)

		// Documents
		r.Get("/documents", docHandler.List)
		r.Post("/documents/upload", docHandler.Upload)
		r.Delete("/documents/{id}", docHandler.Delete)
		r.Get("/documents/{id}/status", docHandler.GetStatus)
		r.Get("/documents/{id}/analysis", analyzerHandler.GetAnalysis)

		// Chat
		r.Get("/chat/{docId}/history", chatHandler.GetHistory)
		r.Delete("/chat/{docId}/history", chatHandler.ClearHistory)
		r.Post("/chat/{docId}/query", chatHandler.Query)

		// Health Tracking
		r.Get("/health/metrics", healthHandler.GetMetrics)
		r.Post("/health/metrics", healthHandler.CreateMetric)
		r.Get("/health/goals", healthHandler.GetGoals)
		r.Post("/health/goals", healthHandler.CreateGoal)
		r.Patch("/health/goals/{id}", healthHandler.UpdateGoal)

		// Meal Plans
		r.Post("/meal-plan/generate", mealPlanHandler.Generate)
		r.Get("/meal-plan/latest", mealPlanHandler.GetLatest)
		r.Delete("/meal-plan/{id}", mealPlanHandler.Delete)

		// Workout Plans
		r.Post("/workout-plan/generate", workoutPlanHandler.Generate)
		r.Get("/workout-plan/latest", workoutPlanHandler.GetLatest)
		r.Delete("/workout-plan/{id}", workoutPlanHandler.Delete)
	})

	// HTTP Server Configuration
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 10 * time.Minute, // High timeout to accommodate query streams
		IdleTimeout:  60 * time.Second,
	}

	// Server execution
	go func() {
		slog.Info("Starting FitMind backend server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ListenAndServe failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("Shutting down backend server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced shutdown", "error", err)
	}

	slog.Info("Backend server gracefully stopped")
}
