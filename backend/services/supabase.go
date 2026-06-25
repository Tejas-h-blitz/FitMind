package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"fitmind/models"
)

type SupabaseService struct {
	url    string
	key    string
	client *http.Client
}

func NewSupabaseService() *SupabaseService {
	return &SupabaseService{
		url:    os.Getenv("SUPABASE_URL"),
		key:    os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *SupabaseService) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, s.url+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", s.key)
	req.Header.Set("Authorization", "Bearer "+s.key)
	return req, nil
}

func (s *SupabaseService) doRequest(req *http.Request) ([]byte, int, error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return respBody, resp.StatusCode, nil
}

// User Profiles
func (s *SupabaseService) GetUserProfile(userID string) (*models.UserProfile, error) {
	req, err := s.newRequest("GET", "/rest/v1/user_profiles?id=eq."+userID+"&select=*", nil)
	if err != nil {
		return nil, err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to get user profile, status: %d, response: %s", code, string(body))
	}

	var profiles []models.UserProfile
	if err := json.Unmarshal(body, &profiles); err != nil {
		return nil, err
	}

	if len(profiles) == 0 {
		return nil, nil
	}
	return &profiles[0], nil
}

func (s *SupabaseService) UpdateUserProfile(userID string, updates map[string]interface{}) error {
	jsonBytes, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	req, err := s.newRequest("PATCH", "/rest/v1/user_profiles?id=eq."+userID, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 200 && code != 204 {
		return fmt.Errorf("failed to update user profile, status: %d, response: %s", code, string(body))
	}
	return nil
}

func (s *SupabaseService) CreateUserProfile(profile *models.UserProfile) error {
	jsonBytes, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	req, err := s.newRequest("POST", "/rest/v1/user_profiles", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 201 && code != 200 {
		return fmt.Errorf("failed to create user profile, status: %d, response: %s", code, string(body))
	}
	return nil
}

// Documents
func (s *SupabaseService) CreateDocument(doc *models.Document) error {
	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req, err := s.newRequest("POST", "/rest/v1/documents", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 201 && code != 200 {
		return fmt.Errorf("failed to create document, status: %d, response: %s", code, string(body))
	}
	return nil
}

func (s *SupabaseService) GetDocument(docID string) (*models.Document, error) {
	req, err := s.newRequest("GET", "/rest/v1/documents?id=eq."+docID+"&select=*", nil)
	if err != nil {
		return nil, err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to get document, status: %d, response: %s", code, string(body))
	}

	var docs []models.Document
	if err := json.Unmarshal(body, &docs); err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, nil
	}
	return &docs[0], nil
}

func (s *SupabaseService) UpdateDocumentStatus(docID string, status string) error {
	updates := map[string]string{"status": status}
	jsonBytes, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	req, err := s.newRequest("PATCH", "/rest/v1/documents?id=eq."+docID, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 200 && code != 204 {
		return fmt.Errorf("failed to update document status, status: %d, response: %s", code, string(body))
	}
	return nil
}

func (s *SupabaseService) ListDocuments(userID string) ([]models.Document, error) {
	req, err := s.newRequest("GET", "/rest/v1/documents?user_id=eq."+userID+"&order=created_at.desc&select=*", nil)
	if err != nil {
		return nil, err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to list documents, status: %d, response: %s", code, string(body))
	}

	var docs []models.Document
	if err := json.Unmarshal(body, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (s *SupabaseService) DeleteDocument(docID string) error {
	req, err := s.newRequest("DELETE", "/rest/v1/documents?id=eq."+docID, nil)
	if err != nil {
		return err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 200 && code != 204 {
		return fmt.Errorf("failed to delete document, status: %d, response: %s", code, string(body))
	}
	return nil
}

// Chat Messages
func (s *SupabaseService) CreateMessage(msg *models.Message) error {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := s.newRequest("POST", "/rest/v1/messages", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 201 && code != 200 {
		return fmt.Errorf("failed to create message, status: %d, response: %s", code, string(body))
	}
	return nil
}

func (s *SupabaseService) GetChatHistory(docID string, userID string) ([]models.Message, error) {
	req, err := s.newRequest("GET", "/rest/v1/messages?doc_id=eq."+docID+"&user_id=eq."+userID+"&order=created_at.asc&select=*", nil)
	if err != nil {
		return nil, err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to list messages, status: %d, response: %s", code, string(body))
	}

	var messages []models.Message
	if err := json.Unmarshal(body, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *SupabaseService) ClearChatHistory(docID string, userID string) error {
	req, err := s.newRequest("DELETE", "/rest/v1/messages?doc_id=eq."+docID+"&user_id=eq."+userID, nil)
	if err != nil {
		return err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 200 && code != 204 {
		return fmt.Errorf("failed to clear messages, status: %d, response: %s", code, string(body))
	}
	return nil
}

// Health Metrics
func (s *SupabaseService) CreateHealthMetric(metric *models.HealthMetric) error {
	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	req, err := s.newRequest("POST", "/rest/v1/health_metrics", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 201 && code != 200 {
		return fmt.Errorf("failed to create health metric, status: %d, response: %s", code, string(body))
	}
	return nil
}

func (s *SupabaseService) GetHealthMetrics(userID string) ([]models.HealthMetric, error) {
	req, err := s.newRequest("GET", "/rest/v1/health_metrics?user_id=eq."+userID+"&order=recorded_at.asc&select=*", nil)
	if err != nil {
		return nil, err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to list health metrics, status: %d, response: %s", code, string(body))
	}

	var metrics []models.HealthMetric
	if err := json.Unmarshal(body, &metrics); err != nil {
		return nil, err
	}
	return metrics, nil
}

// Goals
func (s *SupabaseService) CreateGoal(goal *models.Goal) error {
	jsonBytes, err := json.Marshal(goal)
	if err != nil {
		return err
	}

	req, err := s.newRequest("POST", "/rest/v1/goals", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 201 && code != 200 {
		return fmt.Errorf("failed to create goal, status: %d, response: %s", code, string(body))
	}
	return nil
}

func (s *SupabaseService) GetGoals(userID string) ([]models.Goal, error) {
	req, err := s.newRequest("GET", "/rest/v1/goals?user_id=eq."+userID+"&order=created_at.desc&select=*", nil)
	if err != nil {
		return nil, err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to list goals, status: %d, response: %s", code, string(body))
	}

	var goals []models.Goal
	if err := json.Unmarshal(body, &goals); err != nil {
		return nil, err
	}
	return goals, nil
}

func (s *SupabaseService) UpdateGoalStatus(goalID string, status string) error {
	updates := map[string]string{"status": status}
	jsonBytes, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	req, err := s.newRequest("PATCH", "/rest/v1/goals?id=eq."+goalID, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 200 && code != 204 {
		return fmt.Errorf("failed to update goal, status: %d, response: %s", code, string(body))
	}
	return nil
}

// Supabase Storage Integration
func (s *SupabaseService) UploadFile(storagePath string, fileData []byte) error {
	// The path in URL should match /storage/v1/object/documents/{storagePath}
	// We create the request using POST
	req, err := s.newRequest("POST", "/storage/v1/object/documents/"+storagePath, bytes.NewReader(fileData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/pdf")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}

	// 200 OK or 201 Created is acceptable
	if code != 200 && code != 201 {
		return fmt.Errorf("failed to upload file, status: %d, response: %s", code, string(body))
	}

	return nil
}

func (s *SupabaseService) DeleteFile(storagePath string) error {
	req, err := s.newRequest("DELETE", "/storage/v1/object/documents/"+storagePath, nil)
	if err != nil {
		return err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}

	if code != 200 && code != 204 {
		return fmt.Errorf("failed to delete file, status: %d, response: %s", code, string(body))
	}

	return nil
}

// Document Analysis Operations
func (s *SupabaseService) GetDocumentAnalysis(docID, userID string) (*models.DocumentAnalysis, error) {
	req, err := s.newRequest("GET", "/rest/v1/document_analyses?doc_id=eq."+docID+"&user_id=eq."+userID+"&select=*", nil)
	if err != nil {
		return nil, err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to get document analysis, status: %d, response: %s", code, string(body))
	}

	var analyses []models.DocumentAnalysis
	if err := json.Unmarshal(body, &analyses); err != nil {
		return nil, err
	}

	if len(analyses) == 0 {
		return nil, nil
	}
	return &analyses[0], nil
}

func (s *SupabaseService) SaveDocumentAnalysis(docID, userID string, analysis *models.DocumentAnalysis) error {
	analysis.DocID = docID
	analysis.UserID = userID
	jsonBytes, err := json.Marshal(analysis)
	if err != nil {
		return err
	}

	req, err := s.newRequest("POST", "/rest/v1/document_analyses", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "resolution=merge-duplicates")

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 201 && code != 200 && code != 204 {
		return fmt.Errorf("failed to save document analysis, status: %d, response: %s", code, string(body))
	}
	return nil
}

// Meal Plan Operations
func (s *SupabaseService) GetLatestMealPlan(userID string) (*models.MealPlan, error) {
	req, err := s.newRequest("GET", "/rest/v1/meal_plans?user_id=eq."+userID+"&order=created_at.desc&limit=1&select=*", nil)
	if err != nil {
		return nil, err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to get latest meal plan, status: %d, response: %s", code, string(body))
	}

	var plans []models.MealPlan
	if err := json.Unmarshal(body, &plans); err != nil {
		return nil, err
	}

	if len(plans) == 0 {
		return nil, nil
	}
	return &plans[0], nil
}

func (s *SupabaseService) DeleteMealPlan(mealPlanID, userID string) error {
	req, err := s.newRequest("DELETE", "/rest/v1/meal_plans?id=eq."+mealPlanID+"&user_id=eq."+userID, nil)
	if err != nil {
		return err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 200 && code != 204 {
		return fmt.Errorf("failed to delete meal plan, status: %d, response: %s", code, string(body))
	}
	return nil
}

// Workout Plan Operations
func (s *SupabaseService) GetLatestWorkoutPlan(userID string) (*models.WorkoutPlan, error) {
	req, err := s.newRequest("GET", "/rest/v1/workout_plans?user_id=eq."+userID+"&order=created_at.desc&limit=1&select=*", nil)
	if err != nil {
		return nil, err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to get latest workout plan, status: %d, response: %s", code, string(body))
	}

	var plans []models.WorkoutPlan
	if err := json.Unmarshal(body, &plans); err != nil {
		return nil, err
	}

	if len(plans) == 0 {
		return nil, nil
	}
	return &plans[0], nil
}

func (s *SupabaseService) DeleteWorkoutPlan(workoutPlanID, userID string) error {
	req, err := s.newRequest("DELETE", "/rest/v1/workout_plans?id=eq."+workoutPlanID+"&user_id=eq."+userID, nil)
	if err != nil {
		return err
	}

	body, code, err := s.doRequest(req)
	if err != nil {
		return err
	}
	if code != 200 && code != 204 {
		return fmt.Errorf("failed to delete workout plan, status: %d, response: %s", code, string(body))
	}
	return nil
}

