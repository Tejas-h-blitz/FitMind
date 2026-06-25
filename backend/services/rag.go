package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type RAGService struct {
	url    string
	client *http.Client
}

func NewRAGService() *RAGService {
	url := os.Getenv("RAG_SERVICE_URL")
	if url == "" {
		url = "http://localhost:8000"
	}
	return &RAGService{
		url:    url,
		// Ingestion and queries can take some time, use a generous timeout
		client: &http.Client{Timeout: 5 * time.Minute},
	}
}

func (r *RAGService) TriggerIngest(userID, docID, storagePath string) error {
	payload := map[string]string{
		"user_id":      userID,
		"doc_id":       docID,
		"storage_path": storagePath,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal ingest payload: %w", err)
	}

	url := fmt.Sprintf("%s/ingest", r.url)
	resp, err := r.client.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to call ingest microservice: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ingestion service returned status code %d", resp.StatusCode)
	}

	return nil
}

func (r *RAGService) QueryStream(query, userID, docID string, chatHistory []interface{}) (*http.Response, error) {
	payload := map[string]interface{}{
		"query":        query,
		"user_id":      userID,
		"doc_id":       docID,
		"chat_history": chatHistory,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query payload: %w", err)
	}

	url := fmt.Sprintf("%s/query", r.url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create query request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do not use the service's default short timeout client. We want to support streaming.
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to query microservice: %w", err)
	}

	return resp, nil
}

func (r *RAGService) GenerateMealPlan(userID, docID, preference string) ([]byte, error) {
	payload := map[string]string{
		"user_id":            userID,
		"doc_id":             docID,
		"dietary_preference": preference,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal meal plan payload: %w", err)
	}

	url := fmt.Sprintf("%s/meal-plan", r.url)
	resp, err := r.client.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to call meal plan microservice: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read meal plan response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meal plan service returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return bodyBytes, nil
}

func (r *RAGService) GenerateWorkoutPlan(userID, fitnessLevel, equipment string, daysPerWeek int) ([]byte, error) {
	payload := map[string]interface{}{
		"user_id":       userID,
		"fitness_level": fitnessLevel,
		"equipment":     equipment,
		"days_per_week": daysPerWeek,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workout plan payload: %w", err)
	}

	url := fmt.Sprintf("%s/workout-plan", r.url)
	resp, err := r.client.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to call workout plan microservice: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read workout plan response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("workout plan service returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return bodyBytes, nil
}

