package models

import "time"

type HealthMetric struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	BMI        float64   `json:"bmi"`
	Height     float64   `json:"height"`
	Weight     float64   `json:"weight"`
	RecordedAt time.Time `json:"recorded_at"`
}

type Goal struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	TargetDate string   `json:"target_date"` // YYYY-MM-DD
	Status    string    `json:"status"`      // active, completed
	CreatedAt time.Time `json:"created_at"`
}
