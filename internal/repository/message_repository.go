package repository

import (
	"fmt"

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
	result := r.db.Model(&model.Message{}).
		Where("id = ?", message.ID).
		Updates(map[string]interface{}{
			"is_sent":    message.IsSent,
			"sent_at":    message.SentAt,
			"message_id": message.MessageID,
			"updated_at": message.UpdatedAt,
		})

	if result.Error != nil {
		return fmt.Errorf("mesaj güncellenemedi: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("mesaj bulunamadı veya güncelleme yapılamadı, ID: %d", message.ID)
	}

	return nil
}

func (r *MessageRepository) GetSentMessages() ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where("is_sent = ?", true).Find(&messages).Error
	return messages, err
}
