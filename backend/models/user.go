package models

import "time"

type UserProfile struct {
	ID          string    `json:"id"`
	FullName    string    `json:"full_name"`
	StreakCount int       `json:"streak_count"`
	LastActive  string    `json:"last_active"` // Format: YYYY-MM-DD
	CreatedAt   time.Time `json:"created_at"`
}
