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

// @Summary Get sent messages
// @Description Get a list of all sent messages
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {array} model.MessageResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/messages/sent [get]
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

// @Summary Get all messages
// @Description Get a list of all messages (both sent and unsent)
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {array} model.MessageResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/messages [get]
func (h *MessageHandler) GetAllMessages(c *gin.Context) {
	messages, err := h.service.GetAllMessages()
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
