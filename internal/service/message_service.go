package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ridvanuyn/messaging-system-go/internal/config"
	"github.com/ridvanuyn/messaging-system-go/internal/domain"
	"github.com/ridvanuyn/messaging-system-go/internal/repository"
)

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

func NewMessageService(repo repository.MessageRepository, config *config.Config) MessageService {
	return &messageService{
		repo:      repo,
		config:    config,
		client:    &http.Client{Timeout: 10 * time.Second},
		batchSize: 2, // As specified in requirements, send 2 messages at a time
	}
}

func (s *messageService) SendMessages(ctx context.Context) error {
	messages, err := s.repo.GetUnsentMessages(ctx, s.batchSize)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

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

func (s *messageService) GetSentMessages(ctx context.Context) ([]domain.Message, error) {
	return s.repo.GetSentMessages(ctx)
}

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


	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Yanıt okuma hatası: %v", err)
		return "", err
	}


	log.Printf("Webhook Response - Status Code: %d, Headers: %v, Body: %s", 
		resp.StatusCode, resp.Header, string(bodyBytes))


	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		messageID := fmt.Sprintf("success-%d", time.Now().UnixNano())
		log.Printf("Webhook Success! - Status Code: %d,  MessageID: %s", 
			resp.StatusCode, messageID)
		return messageID, nil
	} else {
		log.Printf("Webhook Failiure: Status Code: %d, URL: %s, Request: %s",
			resp.StatusCode, s.config.WebhookURL, string(jsonData))
		return "", fmt.Errorf("Webhook returned unexpected status code: %d", resp.StatusCode)
	}
}
