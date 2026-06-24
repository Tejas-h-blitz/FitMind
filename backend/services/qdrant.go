package services

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type QdrantService struct {
	url    string
	client *http.Client
}

func NewQdrantService() *QdrantService {
	url := os.Getenv("QDRANT_URL")
	if url == "" {
		url = "http://localhost:6333"
	}
	return &QdrantService{
		url:    url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (q *QdrantService) DeleteCollection(userID, docID string) error {
	// Replicate collection name cleaning logic from python service
	cleanUser := strings.ReplaceAll(userID, "-", "_")
	cleanDoc := strings.ReplaceAll(docID, "-", "_")
	collectionName := fmt.Sprintf("col_%s_%s", cleanUser, cleanDoc)

	url := fmt.Sprintf("%s/collections/%s", q.url, collectionName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete collection request: %w", err)
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Qdrant: %w", err)
	}
	defer resp.Body.Close()

	// 200 is success. 404 means collection was already deleted or never existed, which is also fine.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code from Qdrant: %d", resp.StatusCode)
	}

	return nil
}
