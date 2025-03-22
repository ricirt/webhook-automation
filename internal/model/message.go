package model

import (
	"time"
)

type Message struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	PhoneNumber string    `gorm:"type:varchar(20);not null" json:"phone_number"`
	IsSent      bool      `gorm:"default:false" json:"is_sent"`
	SentAt      time.Time `gorm:"default:null" json:"sent_at,omitempty"`
	MessageID   string    `gorm:"type:varchar(100);default:null" json:"message_id,omitempty"`
}

type MessageRequest struct {
	Content     string `json:"content" binding:"required,max=160"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type MessageResponse struct {
	ID          uint      `json:"id"`
	Content     string    `json:"content"`
	PhoneNumber string    `json:"phone_number"`
	IsSent      bool      `json:"is_sent"`
	SentAt      time.Time `json:"sent_at,omitempty"`
	MessageID   string    `json:"message_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
