package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"fitmind/middleware"
	"fitmind/models"
	"fitmind/services"
	"github.com/go-chi/chi/v5"
)

type WorkoutPlanHandler struct {
	supabase *services.SupabaseService
	rag      *services.RAGService
}

func NewWorkoutPlanHandler(sb *services.SupabaseService, rg *services.RAGService) *WorkoutPlanHandler {
	return &WorkoutPlanHandler{
		supabase: sb,
		rag:      rg,
	}
}

func (h *WorkoutPlanHandler) Generate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	var reqBody struct {
		FitnessLevel string `json:"fitness_level"`
		Equipment    string `json:"equipment"`
		DaysPerWeek  int    `json:"days_per_week"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.FitnessLevel == "" || reqBody.Equipment == "" || reqBody.DaysPerWeek <= 0 {
		SendJSON(w, http.StatusBadRequest, false, nil, "fitness_level, equipment, and days_per_week are required")
		return
	}

	// Call the Python RAG service to generate and save the plan
	responseBytes, err := h.rag.GenerateWorkoutPlan(userID, reqBody.FitnessLevel, reqBody.Equipment, reqBody.DaysPerWeek)
	if err != nil {
		slog.Error("RAG service failed to generate workout plan", "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to generate workout plan")
		return
	}

	// Python returns {"success": true, "data": models.WorkoutPlan}
	var apiRes struct {
		Success bool               `json:"success"`
		Data    models.WorkoutPlan `json:"data"`
		Error   string             `json:"error"`
	}

	if err := json.Unmarshal(responseBytes, &apiRes); err != nil {
		slog.Error("Failed to parse RAG service workout plan response", "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to process workout plan response")
		return
	}

	if !apiRes.Success {
		SendJSON(w, http.StatusInternalServerError, false, nil, apiRes.Error)
		return
	}

	SendJSON(w, http.StatusOK, true, apiRes.Data, "")
}

func (h *WorkoutPlanHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	plan, err := h.supabase.GetLatestWorkoutPlan(userID)
	if err != nil {
		slog.Error("Failed to fetch latest workout plan", "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to load workout plan")
		return
	}

	if plan == nil {
		SendJSON(w, http.StatusNotFound, false, nil, "No workout plan found")
		return
	}

	SendJSON(w, http.StatusOK, true, plan, "")
}

func (h *WorkoutPlanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	workoutPlanID := chi.URLParam(r, "id")
	if workoutPlanID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Workout Plan ID required")
		return
	}

	err := h.supabase.DeleteWorkoutPlan(workoutPlanID, userID)
	if err != nil {
		slog.Error("Failed to delete workout plan", "workout_plan_id", workoutPlanID, "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to delete workout plan")
		return
	}

	SendJSON(w, http.StatusOK, true, map[string]string{"message": "Workout plan deleted successfully"}, "")
}
