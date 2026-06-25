package handlers

import (
	"log/slog"
	"net/http"
	"fitmind/middleware"
	"fitmind/services"
	"github.com/go-chi/chi/v5"
)

type AnalyzerHandler struct {
	supabase *services.SupabaseService
}

func NewAnalyzerHandler(sb *services.SupabaseService) *AnalyzerHandler {
	return &AnalyzerHandler{supabase: sb}
}

func (h *AnalyzerHandler) GetAnalysis(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	docID := chi.URLParam(r, "id")
	if docID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Document ID required")
		return
	}

	analysis, err := h.supabase.GetDocumentAnalysis(docID, userID)
	if err != nil {
		slog.Error("Failed to fetch document analysis", "doc_id", docID, "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to load document analysis")
		return
	}

	if analysis == nil {
		SendJSON(w, http.StatusNotFound, false, nil, "Analysis not found or still processing")
		return
	}

	SendJSON(w, http.StatusOK, true, analysis, "")
}
