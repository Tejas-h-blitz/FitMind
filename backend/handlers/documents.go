package handlers

import (
	"crypto/rand"
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

type DocumentsHandler struct {
	supabase *services.SupabaseService
	qdrant   *services.QdrantService
	rag      *services.RAGService
}

func NewDocumentsHandler(sb *services.SupabaseService, qd *services.QdrantService, rg *services.RAGService) *DocumentsHandler {
	return &DocumentsHandler{
		supabase: sb,
		qdrant:   qd,
		rag:      rg,
	}
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	// Formats as 8-4-4-4-12 representation
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (h *DocumentsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	docs, err := h.supabase.ListDocuments(userID)
	if err != nil {
		slog.Error("Failed to list documents", "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to list documents")
		return
	}

	SendJSON(w, http.StatusOK, true, docs, "")
}

func (h *DocumentsHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	// Limit upload size to 15MB
	r.Body = http.MaxBytesReader(w, r.Body, 15<<20)
	err := r.ParseMultipartForm(15 << 20)
	if err != nil {
		SendJSON(w, http.StatusBadRequest, false, nil, "File too large or failed to parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		SendJSON(w, http.StatusBadRequest, false, nil, "Missing file in request")
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		SendJSON(w, http.StatusBadRequest, false, nil, "Only PDF files are supported")
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read uploaded file", "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to process file")
		return
	}

	docID := generateUUID()
	storagePath := fmt.Sprintf("uploads/%s/%s.pdf", userID, docID)

	slog.Info("Uploading file to Supabase Storage", "storage_path", storagePath)
	err = h.supabase.UploadFile(storagePath, fileData)
	if err != nil {
		slog.Error("Failed to upload file to storage", "user_id", userID, "doc_id", docID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to upload file to storage")
		return
	}

	doc := &models.Document{
		ID:          docID,
		UserID:      userID,
		Name:        header.Filename,
		Size:        header.Size,
		StoragePath: storagePath,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	slog.Info("Saving document metadata to DB", "doc_id", docID)
	err = h.supabase.CreateDocument(doc)
	if err != nil {
		slog.Error("Failed to save document metadata", "doc_id", docID, "error", err)
		// Clean up the uploaded file from storage
		_ = h.supabase.DeleteFile(storagePath)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to save document metadata")
		return
	}

	// Trigger RAG Ingestion asynchronously in the background
	go func(uID, dID, path string) {
		slog.Info("Starting background ingestion", "doc_id", dID)
		
		err := h.supabase.UpdateDocumentStatus(dID, "processing")
		if err != nil {
			slog.Error("Failed to update status to processing", "doc_id", dID, "error", err)
		}

		err = h.rag.TriggerIngest(uID, dID, path)
		if err != nil {
			slog.Error("Background ingestion failed", "doc_id", dID, "error", err)
			_ = h.supabase.UpdateDocumentStatus(dID, "failed")
			return
		}

		slog.Info("Background ingestion succeeded", "doc_id", dID)
		_ = h.supabase.UpdateDocumentStatus(dID, "ready")
	}(userID, docID, storagePath)

	SendJSON(w, http.StatusAccepted, true, doc, "")
}

func (h *DocumentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	// Verify document belongs to the user
	doc, err := h.supabase.GetDocument(docID)
	if err != nil {
		slog.Error("Failed to fetch document to delete", "doc_id", docID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to verify document ownership")
		return
	}
	if doc == nil || doc.UserID != userID {
		SendJSON(w, http.StatusNotFound, false, nil, "Document not found")
		return
	}

	// 1. Delete Qdrant collection
	slog.Info("Dropping Qdrant collection for document", "doc_id", docID)
	err = h.qdrant.DeleteCollection(userID, docID)
	if err != nil {
		slog.Error("Failed to delete Qdrant collection", "doc_id", docID, "error", err)
		// We log and continue so we clean up storage and postgres even if Qdrant fails
	}

	// 2. Delete Supabase Storage File
	slog.Info("Deleting PDF from storage", "storage_path", doc.StoragePath)
	err = h.supabase.DeleteFile(doc.StoragePath)
	if err != nil {
		slog.Error("Failed to delete file from Storage", "storage_path", doc.StoragePath, "error", err)
		// Continue to clean up database metadata
	}

	// 3. Delete database record
	slog.Info("Deleting database record", "doc_id", docID)
	err = h.supabase.DeleteDocument(docID)
	if err != nil {
		slog.Error("Failed to delete document from database", "doc_id", docID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to delete document metadata")
		return
	}

	SendJSON(w, http.StatusOK, true, map[string]string{"message": "Document deleted successfully"}, "")
}

func (h *DocumentsHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
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

	doc, err := h.supabase.GetDocument(docID)
	if err != nil {
		slog.Error("Failed to get document status", "doc_id", docID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to fetch document status")
		return
	}

	if doc == nil || doc.UserID != userID {
		SendJSON(w, http.StatusNotFound, false, nil, "Document not found")
		return
	}

	SendJSON(w, http.StatusOK, true, map[string]string{
		"id":     doc.ID,
		"status": doc.Status,
	}, "")
}
