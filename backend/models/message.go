package models

import "time"

type SourceChunk struct {
	Text       string `json:"text"`
	PageNumber int    `json:"page_number"`
}

type Message struct {
	ID        string        `json:"id"`
	DocID     string        `json:"doc_id"`
	UserID    string        `json:"user_id"`
	Role      string        `json:"role"` // user/assistant
	Content   string        `json:"content"`
	Sources   []SourceChunk `json:"sources"` // Serialized to jsonb in Supabase
	CreatedAt time.Time     `json:"created_at"`
}
