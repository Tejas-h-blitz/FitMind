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

type MealPlanHandler struct {
	supabase *services.SupabaseService
	rag      *services.RAGService
}

func NewMealPlanHandler(sb *services.SupabaseService, rg *services.RAGService) *MealPlanHandler {
	return &MealPlanHandler{
		supabase: sb,
		rag:      rg,
	}
}

func (h *MealPlanHandler) Generate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	var reqBody struct {
		DocID             string `json:"doc_id"`
		DietaryPreference string `json:"dietary_preference"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.DocID == "" || reqBody.DietaryPreference == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "doc_id and dietary_preference are required")
		return
	}

	// Call the Python RAG service to generate and save the plan
	responseBytes, err := h.rag.GenerateMealPlan(userID, reqBody.DocID, reqBody.DietaryPreference)
	if err != nil {
		slog.Error("RAG service failed to generate meal plan", "user_id", userID, "doc_id", reqBody.DocID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to generate meal plan")
		return
	}

	// Python returns {"success": true, "data": models.MealPlan}
	var apiRes struct {
		Success bool            `json:"success"`
		Data    models.MealPlan `json:"data"`
		Error   string          `json:"error"`
	}

	if err := json.Unmarshal(responseBytes, &apiRes); err != nil {
		slog.Error("Failed to parse RAG service meal plan response", "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to process meal plan response")
		return
	}

	if !apiRes.Success {
		SendJSON(w, http.StatusInternalServerError, false, nil, apiRes.Error)
		return
	}

	SendJSON(w, http.StatusOK, true, apiRes.Data, "")
}

func (h *MealPlanHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	plan, err := h.supabase.GetLatestMealPlan(userID)
	if err != nil {
		slog.Error("Failed to fetch latest meal plan", "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to load meal plan")
		return
	}

	if plan == nil {
		SendJSON(w, http.StatusNotFound, false, nil, "No meal plan found")
		return
	}

	SendJSON(w, http.StatusOK, true, plan, "")
}

func (h *MealPlanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	mealPlanID := chi.URLParam(r, "id")
	if mealPlanID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Meal Plan ID required")
		return
	}

	err := h.supabase.DeleteMealPlan(mealPlanID, userID)
	if err != nil {
		slog.Error("Failed to delete meal plan", "meal_plan_id", mealPlanID, "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to delete meal plan")
		return
	}

	SendJSON(w, http.StatusOK, true, map[string]string{"message": "Meal plan deleted successfully"}, "")
}
