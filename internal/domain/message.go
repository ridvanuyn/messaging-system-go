package domain

import (
	"time"
)

// Message represents a message in the database
type Message struct {
	ID        int64     `json:"id" db:"id"`
	To        string    `json:"to" db:"to"`
	Content   string    `json:"content" db:"content"`
	Sent      bool      `json:"sent" db:"sent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	SentAt    time.Time `json:"sent_at,omitempty" db:"sent_at"`
	MessageID string    `json:"message_id,omitempty" db:"message_id"`
}

// MessageResponse represents the response from the webhook site
type MessageResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}
