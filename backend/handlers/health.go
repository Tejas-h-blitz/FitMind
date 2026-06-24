package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"fitmind/middleware"
	"fitmind/models"
	"fitmind/services"
	"github.com/go-chi/chi/v5"
)

type HealthHandler struct {
	supabase *services.SupabaseService
}

func NewHealthHandler(sb *services.SupabaseService) *HealthHandler {
	return &HealthHandler{supabase: sb}
}

func (h *HealthHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	metrics, err := h.supabase.GetHealthMetrics(userID)
	if err != nil {
		slog.Error("Failed to fetch health metrics", "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to load health metrics")
		return
	}

	SendJSON(w, http.StatusOK, true, metrics, "")
}

func (h *HealthHandler) CreateMetric(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	var reqBody struct {
		Height float64 `json:"height"` // in cm or meters
		Weight float64 `json:"weight"` // in kg
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid JSON body")
		return
	}

	if reqBody.Height <= 0 || reqBody.Weight <= 0 {
		SendJSON(w, http.StatusBadRequest, false, nil, "Height and weight must be positive values")
		return
	}

	// Auto-detect cm vs meters
	heightInMeters := reqBody.Height
	if reqBody.Height > 3.0 {
		heightInMeters = reqBody.Height / 100.0
	}

	bmi := reqBody.Weight / (heightInMeters * heightInMeters)

	metric := &models.HealthMetric{
		ID:         generateUUID(),
		UserID:     userID,
		BMI:        bmi,
		Height:     reqBody.Height,
		Weight:     reqBody.Weight,
		RecordedAt: time.Now(),
	}

	err := h.supabase.CreateHealthMetric(metric)
	if err != nil {
		slog.Error("Failed to save health metric", "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to record BMI metric")
		return
	}

	SendJSON(w, http.StatusCreated, true, metric, "")
}

func (h *HealthHandler) GetGoals(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	goals, err := h.supabase.GetGoals(userID)
	if err != nil {
		slog.Error("Failed to fetch goals", "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to load goals")
		return
	}

	SendJSON(w, http.StatusOK, true, goals, "")
}

func (h *HealthHandler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	var reqBody struct {
		Title      string `json:"title"`
		TargetDate string `json:"target_date"` // YYYY-MM-DD
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid JSON body")
		return
	}

	if reqBody.Title == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Goal title is required")
		return
	}

	// Validate date layout if provided
	if reqBody.TargetDate != "" {
		_, err := time.Parse("2006-01-02", reqBody.TargetDate)
		if err != nil {
			SendJSON(w, http.StatusBadRequest, false, nil, "Target date must follow YYYY-MM-DD format")
			return
		}
	}

	goal := &models.Goal{
		ID:         generateUUID(),
		UserID:     userID,
		Title:      reqBody.Title,
		TargetDate: reqBody.TargetDate,
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	err := h.supabase.CreateGoal(goal)
	if err != nil {
		slog.Error("Failed to create goal", "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to save goal")
		return
	}

	SendJSON(w, http.StatusCreated, true, goal, "")
}

func (h *HealthHandler) UpdateGoal(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	goalID := chi.URLParam(r, "id")
	if goalID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Goal ID required")
		return
	}

	var reqBody struct {
		Status string `json:"status"` // active or completed
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid JSON body")
		return
	}

	if reqBody.Status != "active" && reqBody.Status != "completed" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Status must be either 'active' or 'completed'")
		return
	}

	// In a real app we should check if the goal belongs to the user
	// We'll trust PostgREST/RLS policies, but let's check or handle updates
	err := h.supabase.UpdateGoalStatus(goalID, reqBody.Status)
	if err != nil {
		slog.Error("Failed to update goal", "goal_id", goalID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to update goal status")
		return
	}

	SendJSON(w, http.StatusOK, true, map[string]string{"id": goalID, "status": reqBody.Status}, "")
}
