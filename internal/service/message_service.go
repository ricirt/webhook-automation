package service

import (
	"github.com/ricirt/webhook-automation/internal/model"
	"github.com/ricirt/webhook-automation/internal/repository"
)

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

type MessageService struct {
	repo        *repository.MessageRepository
	isRunning   bool
	stopChannel chan struct{}
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{
		repo:        repo,
		stopChannel: make(chan struct{}),
	}
}

func (s *MessageService) GetSentMessages() ([]model.Message, error) {
	return s.repo.GetSentMessages()
}

func (s *MessageService) GetAllMessages() ([]model.Message, error) {
	return s.repo.GetAllMessages()
}
