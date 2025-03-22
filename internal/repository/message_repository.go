package repository

import (
	"github.com/ricirt/webhook-automation/internal/model"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) GetUnsentMessages(limit int) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where("is_sent = ?", false).Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) UpdateMessage(message *model.Message) error {
	result := r.db.Model(message).Updates(map[string]interface{}{
		"is_sent":    true,
		"sent_at":    message.SentAt,
		"message_id": message.MessageID,
	})
	return result.Error
}

func (r *MessageRepository) GetSentMessages() ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where("is_sent = ?", true).Find(&messages).Error
	return messages, err
}
