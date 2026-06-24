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
