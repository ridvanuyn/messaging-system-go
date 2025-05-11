package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridvanuyn/messaging-system-go/internal/service"
	"github.com/ridvanuyn/messaging-system-go/internal/worker"
)

// Handler contains API handlers
type Handler struct {
	messageService service.MessageService
	scheduler      *worker.Scheduler
}

// NewHandler creates a new handler
func NewHandler(messageService service.MessageService, scheduler *worker.Scheduler) *Handler {
	return &Handler{
		messageService: messageService,
		scheduler:      scheduler,
	}
}

// StartScheduler starts the message scheduler
// @Summary Start message scheduler
// @Description Start automatic message sending
// @Tags scheduler
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /scheduler/start [post]
func (h *Handler) StartScheduler(c *gin.Context) {
	if h.scheduler.Start() {
		c.JSON(http.StatusOK, gin.H{"status": "Message scheduler started"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message scheduler already running"})
	}
}

// StopScheduler stops the message scheduler
// @Summary Stop message scheduler
// @Description Stop automatic message sending
// @Tags scheduler
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /scheduler/stop [post]
func (h *Handler) StopScheduler(c *gin.Context) {
	if h.scheduler.Stop() {
		c.JSON(http.StatusOK, gin.H{"status": "Message scheduler stopped"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message scheduler already stopped"})
	}
}

// GetSentMessages gets sent messages
// @Summary List sent messages
// @Description Get all sent messages from database
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {array} domain.Message
// @Router /messages [get]
func (h *Handler) GetSentMessages(c *gin.Context) {
	messages, err := h.messageService.GetSentMessages(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetSchedulerStatus gets scheduler status
// @Summary Check scheduler status
// @Description Check if scheduler is currently running
// @Tags scheduler
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /scheduler/status [get]
func (h *Handler) GetSchedulerStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"running": h.scheduler.IsRunning()})
}
