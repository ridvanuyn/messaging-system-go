package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ridvanuyn/messaging-system-gointernal/config"
	"github.com/ridvanuyn/messaging-system-gointernal/domain"
	"github.com/ridvanuyn/messaging-system-gointernal/repository"
)

// MessageService defines the message operations
type MessageService interface {
	SendMessages(ctx context.Context) error
	GetSentMessages(ctx context.Context) ([]domain.Message, error)
}

type messageService struct {
	repo      repository.MessageRepository
	config    *config.Config
	client    *http.Client
	batchSize int
}

// NewMessageService creates a new message service instance
func NewMessageService(repo repository.MessageRepository, config *config.Config) MessageService {
	return &messageService{
		repo:      repo,
		config:    config,
		client:    &http.Client{Timeout: 10 * time.Second},
		batchSize: 2, // As specified in requirements, send 2 messages at a time
	}
}

// SendMessages sends unsent messages
func (s *messageService) SendMessages(ctx context.Context) error {
	// Get unsent messages from database
	messages, err := s.repo.GetUnsentMessages(ctx, s.batchSize)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	// For each message
	for _, msg := range messages {
		// Check content length
		if len(msg.Content) > s.config.MaxContentLength {
			log.Printf("Message content too long: %d, max: %d", len(msg.Content), s.config.MaxContentLength)
			continue
		}

		// Send to webhook
		messageID, err := s.sendToWebhook(ctx, msg)
		if err != nil {
			log.Printf("Failed to send message ID: %d, error: %v", msg.ID, err)
			continue
		}

		// Mark as sent in database
		err = s.repo.MarkAsSent(ctx, msg.ID, messageID)
		if err != nil {
			log.Printf("Failed to mark message as sent ID: %d, error: %v", msg.ID, err)
			continue
		}

		// Bonus: Cache to Redis
		sentTime := time.Now()
		err = s.repo.CacheMessageID(ctx, messageID, sentTime)
		if err != nil {
			log.Printf("Failed to cache message ID: %s, error: %v", messageID, err)
		}
	}

	return nil
}

// GetSentMessages gets sent messages
func (s *messageService) GetSentMessages(ctx context.Context) ([]domain.Message, error) {
	return s.repo.GetSentMessages(ctx)
}

// sendToWebhook sends a message to the webhook
func (s *messageService) sendToWebhook(ctx context.Context, msg domain.Message) (string, error) {
	payload := struct {
		To      string `json:"to"`
		Content string `json:"content"`
	}{
		To:      msg.To,
		Content: msg.Content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", s.config.AuthKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for 202 Accepted
	if resp.StatusCode != http.StatusAccepted {
		return "", errors.New("webhook returned unexpected status code")
	}

	// Parse response
	var response domain.MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.MessageID, nil
}
