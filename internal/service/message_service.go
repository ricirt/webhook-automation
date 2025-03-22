package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ricirt/webhook-automation/internal/config"
	"github.com/ricirt/webhook-automation/internal/model"
	"github.com/ricirt/webhook-automation/internal/repository"

	"github.com/go-redis/redis/v8"
)

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

type MessageService struct {
	repo        *repository.MessageRepository
	config      *config.Config
	redis       *redis.Client
	isRunning   bool
	stopChannel chan struct{}
}

func NewMessageService(repo *repository.MessageRepository, config *config.Config, redis *redis.Client) *MessageService {
	return &MessageService{
		repo:        repo,
		config:      config,
		redis:       redis,
		isRunning:   false,
		stopChannel: make(chan struct{}),
	}
}

func (s *MessageService) StartSending() error {
	if s.isRunning {
		return fmt.Errorf("mesaj gönderme servisi zaten çalışıyor")
	}

	s.isRunning = true
	s.stopChannel = make(chan struct{}) // Yeni kanal oluştur
	go s.sendMessagesLoop()
	fmt.Println("Mesaj gönderme servisi başlatıldı")
	return nil
}

func (s *MessageService) StopSending() error {
	if !s.isRunning {
		return fmt.Errorf("mesaj gönderme servisi zaten durdurulmuş")
	}

	s.isRunning = false
	close(s.stopChannel)
	fmt.Println("Mesaj gönderme servisi durdurma sinyali gönderildi")
	return nil
}

func (s *MessageService) sendMessagesLoop() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	if err := s.sendMessages(); err != nil {
		fmt.Printf("Mesaj gönderiminde hata: %v\n", err)
	}

	for {
		select {
		case <-s.stopChannel:
			return
		case <-ticker.C:
			if !s.isRunning {
				return
			}
			if err := s.sendMessages(); err != nil {
				fmt.Printf("Mesaj gönderiminde hata: %v\n", err)
			}
		}
	}
}

func (s *MessageService) sendMessages() error {
	messages, err := s.repo.GetUnsentMessages(2)
	if err != nil {
		return fmt.Errorf("gönderilmemiş mesajlar alınamadı: %v", err)
	}

	if len(messages) == 0 {
		return nil
	}

	for i := 0; i < len(messages); i++ {
		message := &messages[i]
		fmt.Printf("Mesaj işleniyor - ID: %d\n", message.ID)

		if err := s.sendMessage(message); err != nil {
			fmt.Printf("Mesaj gönderilemedi - ID: %d, Hata: %v\n", message.ID, err)
			continue
		}

		if err := s.repo.UpdateMessage(message); err != nil {
			fmt.Printf("Veritabanı güncellenemedi - ID: %d, Hata: %v\n", message.ID, err)
			continue
		}
	}

	return nil
}

func (s *MessageService) sendMessage(message *model.Message) error {
	payload := map[string]string{
		"phoneNumber": message.PhoneNumber,
		"content":     message.Content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("payload oluşturulamadı: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", s.config.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("istek oluşturulamadı: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("istek gönderilemedi: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("yanıt okunamadı: %v", err)
	}

	var webhookResp WebhookResponse
	if err := json.Unmarshal(body, &webhookResp); err != nil {
		return fmt.Errorf("webhook yanıtı ayrıştırılamadı: %v", err)
	}

	if webhookResp.MessageID == "" {
		return fmt.Errorf("messageId boş")
	}

	now := time.Now()
	message.MessageID = webhookResp.MessageID
	message.IsSent = true
	message.SentAt = now
	message.UpdatedAt = now

	if err := s.cacheMessageDetails(message.ID, webhookResp.MessageID, now); err != nil {
		fmt.Printf("Redis'e kaydedilemedi - ID: %d, Hata: %v\n", message.ID, err)
	}

	return nil
}

func (s *MessageService) GetSentMessages() ([]model.Message, error) {
	return s.repo.GetSentMessages()
}

func (s *MessageService) cacheMessageDetails(messageID uint, messageResponseID string, sentAt time.Time) error {
	ctx := context.Background()
	key := fmt.Sprintf("message:%d", messageID)
	value := map[string]interface{}{
		"message_id": messageResponseID,
		"sent_at":    sentAt.Format(time.RFC3339),
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache değeri oluşturulamadı: %v", err)
	}

	if err := s.redis.Set(ctx, key, jsonValue, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("redis'e kaydedilemedi: %v", err)
	}

	return nil
}
