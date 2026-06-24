package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"fitmind/middleware"
	"fitmind/models"
	"fitmind/services"
	"github.com/go-chi/chi/v5"
)

type ChatHandler struct {
	supabase *services.SupabaseService
	rag      *services.RAGService
}

func NewChatHandler(sb *services.SupabaseService, rg *services.RAGService) *ChatHandler {
	return &ChatHandler{
		supabase: sb,
		rag:      rg,
	}
}

func (h *ChatHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	docID := chi.URLParam(r, "docId")
	if docID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Document ID required")
		return
	}

	// Verify document ownership
	doc, err := h.supabase.GetDocument(docID)
	if err != nil || doc == nil || doc.UserID != userID {
		SendJSON(w, http.StatusNotFound, false, nil, "Document not found")
		return
	}

	messages, err := h.supabase.GetChatHistory(docID, userID)
	if err != nil {
		slog.Error("Failed to fetch chat history", "doc_id", docID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to load chat history")
		return
	}

	SendJSON(w, http.StatusOK, true, messages, "")
}

func (h *ChatHandler) ClearHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	docID := chi.URLParam(r, "docId")
	if docID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Document ID required")
		return
	}

	// Verify document ownership
	doc, err := h.supabase.GetDocument(docID)
	if err != nil || doc == nil || doc.UserID != userID {
		SendJSON(w, http.StatusNotFound, false, nil, "Document not found")
		return
	}

	err = h.supabase.ClearChatHistory(docID, userID)
	if err != nil {
		slog.Error("Failed to clear chat history", "doc_id", docID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to clear chat history")
		return
	}

	SendJSON(w, http.StatusOK, true, map[string]string{"message": "History cleared"}, "")
}

func (h *ChatHandler) Query(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	docID := chi.URLParam(r, "docId")
	if docID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Document ID required")
		return
	}

	// Verify document ownership
	doc, err := h.supabase.GetDocument(docID)
	if err != nil || doc == nil || doc.UserID != userID {
		SendJSON(w, http.StatusNotFound, false, nil, "Document not found")
		return
	}

	if doc.Status != "ready" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Document is not ready for chat yet (status: "+doc.Status+")")
		return
	}

	// Parse query
	var reqBody struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.Query == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Query is required")
		return
	}

	// 1. Fetch existing chat history from Postgres
	historyMessages, err := h.supabase.GetChatHistory(docID, userID)
	if err != nil {
		slog.Error("Failed to get chat history for context", "doc_id", docID, "error", err)
		// We can proceed with empty history if fetch fails, but let's log it
	}

	// Format history for RAG microservice
	var chatHistory []interface{}
	for _, m := range historyMessages {
		chatHistory = append(chatHistory, map[string]string{
			"role":    m.Role,
			"content": m.Content,
		})
	}

	// 2. Save User's query to DB
	userMsg := &models.Message{
		ID:        generateUUID(),
		DocID:     docID,
		UserID:    userID,
		Role:      "user",
		Content:   reqBody.Query,
		Sources:   []models.SourceChunk{},
		CreatedAt: time.Now(),
	}
	err = h.supabase.CreateMessage(userMsg)
	if err != nil {
		slog.Error("Failed to save user message", "doc_id", docID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to save message")
		return
	}

	// 3. Request SSE stream from RAG microservice
	slog.Info("Requesting stream from RAG microservice", "doc_id", docID, "query", reqBody.Query)
	ragResp, err := h.rag.QueryStream(reqBody.Query, userID, docID, chatHistory)
	if err != nil {
		slog.Error("Failed to get stream from RAG microservice", "doc_id", docID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to communicate with RAG service: "+err.Error())
		return
	}
	defer ragResp.Body.Close()

	// 4. Set Headers for SSE Response to the SvelteKit frontend
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("Flusher interface not supported for streaming")
		return
	}

	// Stream and accumulate response
	var assistantText strings.Builder
	var sources []models.SourceChunk

	reader := bufio.NewReader(ragResp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			slog.Error("Error reading stream from RAG service", "error", err)
			break
		}

		// Forward the raw SSE line directly to SvelteKit
		_, _ = fmt.Fprint(w, line)
		flusher.Flush()

		// Heuristic: accumulate content and sources from the SSE lines
		// SSE format:
		// data: [{"text": "...", "page_number": 1}]
		// data: "token"
		if strings.HasPrefix(line, "data: ") {
			dataStr := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
			
			if strings.HasPrefix(dataStr, "[") {
				var parsedSources []models.SourceChunk
				if err := json.Unmarshal([]byte(dataStr), &parsedSources); err == nil {
					sources = parsedSources
				}
			} else if strings.HasPrefix(dataStr, "\"") {
				var token string
				if err := json.Unmarshal([]byte(dataStr), &token); err == nil {
					assistantText.WriteString(token)
				}
			}
		}
	}

	// 5. Persist assistant's accumulated response to database
	assistantMsg := &models.Message{
		ID:        generateUUID(),
		DocID:     docID,
		UserID:    userID,
		Role:      "assistant",
		Content:   assistantText.String(),
		Sources:   sources,
		CreatedAt: time.Now(),
	}

	slog.Info("Saving assistant response to database", "doc_id", docID)
	err = h.supabase.CreateMessage(assistantMsg)
	if err != nil {
		slog.Error("Failed to save assistant response message", "doc_id", docID, "error", err)
	}
}
