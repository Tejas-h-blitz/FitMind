package models

import "time"

type Document struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	StoragePath string    `json:"storage_path"`
	Status      string    `json:"status"` // pending, processing, ready, failed
	CreatedAt   time.Time `json:"created_at"`
}

type DocumentHealthMetric struct {
	Name           string  `json:"name"`
	Value          float64 `json:"value"`
	Unit           string  `json:"unit"`
	Status         string  `json:"status"` // normal/low/high/borderline
	ReferenceRange string  `json:"reference_range"`
	PlainEnglish   string  `json:"plain_english"`
}

type DocumentAnalysis struct {
	ID            string                 `json:"id"`
	DocID         string                 `json:"doc_id"`
	UserID        string                 `json:"user_id"`
	Metrics       []DocumentHealthMetric `json:"metrics"`
	Summary       string                 `json:"summary"`
	OverallStatus string                 `json:"overall_status"`
	CreatedAt     time.Time              `json:"created_at"`
}

