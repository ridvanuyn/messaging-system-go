package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ridvanuyn/messaging-system-go/internal/service"
)

// Scheduler manages periodic message sending
type Scheduler struct {
	messageService service.MessageService
	interval       time.Duration
	stopCh         chan struct{}
	isRunning      bool
	mu             sync.Mutex
}

// NewScheduler creates a new scheduler
func NewScheduler(messageService service.MessageService) *Scheduler {
	return &Scheduler{
		messageService: messageService,
		interval:       2 * time.Minute, // As specified in requirements, every 2 minutes
		stopCh:         make(chan struct{}),
		isRunning:      false,
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return false
	}

	s.isRunning = true
	s.stopCh = make(chan struct{})

	// Initial run immediately
	go func() {
		log.Println("Message scheduler started")
		
		// Send messages immediately on start
		s.sendMessages()

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.sendMessages()
			case <-s.stopCh:
				log.Println("Message scheduler stopped")
				return
			}
		}
	}()

	return true
}

// Stop stops the scheduler
func (s *Scheduler) Stop() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return false
	}

	close(s.stopCh)
	s.isRunning = false
	return true
}

// IsRunning checks if the scheduler is running
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isRunning
}

// sendMessages triggers message sending
func (s *Scheduler) sendMessages() {
	log.Println("Sending messages...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.messageService.SendMessages(ctx); err != nil {
		log.Printf("Error during message sending: %v", err)
	}
}
