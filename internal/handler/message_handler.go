package handler

import (
	"net/http"

	"github.com/ricirt/webhook-automation/internal/model"
	"github.com/ricirt/webhook-automation/internal/service"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	service *service.MessageService
}

func NewMessageHandler(service *service.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

// @Summary Start automatic message sending
// @Description Start the automatic message sending process
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /messages/start [post]
func (h *MessageHandler) StartSending(c *gin.Context) {
	if err := h.service.StartSending(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sending started"})
}

// @Summary Stop automatic message sending
// @Description Stop the automatic message sending process
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /messages/stop [post]
func (h *MessageHandler) StopSending(c *gin.Context) {
	if err := h.service.StopSending(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sending stopped"})
}

// @Summary Get sent messages
// @Description Get a list of all sent messages
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {array} model.MessageResponse
// @Failure 500 {object} map[string]string
// @Router /messages/sent [get]
func (h *MessageHandler) GetSentMessages(c *gin.Context) {
	messages, err := h.service.GetSentMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []model.MessageResponse
	for _, msg := range messages {
		response = append(response, model.MessageResponse{
			ID:          msg.ID,
			Content:     msg.Content,
			PhoneNumber: msg.PhoneNumber,
			IsSent:      msg.IsSent,
			SentAt:      msg.SentAt,
			MessageID:   msg.MessageID,
			CreatedAt:   msg.CreatedAt,
			UpdatedAt:   msg.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
