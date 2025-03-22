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

func (r *MessageRepository) GetSentMessages() ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where("is_sent = ?", true).Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) GetAllMessages() ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Find(&messages).Error
	return messages, err
}
