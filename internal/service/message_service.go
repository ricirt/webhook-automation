package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ricirt/webhook-automation/internal/config"
	"github.com/ricirt/webhook-automation/internal/model"
	"github.com/ricirt/webhook-automation/internal/repository"
)

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

type MessageService struct {
	repo        *repository.MessageRepository
	config      *config.Config
	isRunning   bool
	stopChannel chan struct{}
}

func NewMessageService(repo *repository.MessageRepository, config *config.Config) *MessageService {
	return &MessageService{
		repo:        repo,
		config:      config,
		isRunning:   false,
		stopChannel: make(chan struct{}),
	}
}

func (s *MessageService) StartSending() error {
	if s.isRunning {
		return fmt.Errorf("message sending is already running")
	}

	s.isRunning = true
	go s.sendMessagesLoop()
	return nil
}

func (s *MessageService) StopSending() error {
	if !s.isRunning {
		return fmt.Errorf("message sending is not running")
	}

	close(s.stopChannel)
	s.isRunning = false
	return nil
}

func (s *MessageService) sendMessagesLoop() {
	if err := s.sendMessages(); err != nil {
		fmt.Printf("Error in initial message sending: %v\n", err)
	}

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChannel:
			return
		case <-ticker.C:
			if err := s.sendMessages(); err != nil {
				fmt.Printf("Error sending messages: %v\n", err)
			}
		}
	}
}

func (s *MessageService) sendMessages() error {
	messages, err := s.repo.GetUnsentMessages(2)
	if err != nil {
		return fmt.Errorf("failed to get unsent messages: %v", err)
	}

	for _, message := range messages {
		if err := s.sendMessage(&message); err != nil {
			continue
		}

		message.IsSent = true
		message.SentAt = time.Now()

		if err := s.repo.UpdateMessage(&message); err != nil {
			continue
		}
	}

	return nil
}

func (s *MessageService) sendMessage(message *model.Message) error {
	payload := map[string]string{
		"to":      message.PhoneNumber,
		"content": message.Content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", s.config.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var webhookResp WebhookResponse
	if err := json.Unmarshal(body, &webhookResp); err != nil {
		message.MessageID = fmt.Sprintf("msg_%d_%d", message.ID, time.Now().Unix())
		return nil
	}

	message.MessageID = webhookResp.MessageID
	return nil
}

func (s *MessageService) GetSentMessages() ([]model.Message, error) {
	return s.repo.GetSentMessages()
}
