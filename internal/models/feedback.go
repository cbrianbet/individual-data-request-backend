package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Feedback struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RequestID uuid.UUID `gorm:"type:uuid;not null"`
	SenderID  uuid.UUID `gorm:"type:uuid;not null"`
	Message   string    `gorm:"size:2000;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`

	//foreign keys
	Request Requests `gorm:"foreignKey:RequestID"`
}

func NewMessage(db *gorm.DB, requestID, senderID uuid.UUID, message string) (*Feedback, error) {
	chatMessage := &Feedback{
		RequestID: requestID,
		SenderID:  senderID,
		Message:   message,
	}

	if err := db.Create(chatMessage).Error; err != nil {
		return nil, err
	}

	return chatMessage, nil
}

func GetChatMessagesByRequestID(db *gorm.DB, requestID uuid.UUID) ([]Feedback, error) {
	var chatMessages []Feedback
	if err := db.Where("request_id = ?", requestID).Order("created_at asc").Find(&chatMessages).Error; err != nil {
		return nil, err
	}

	return chatMessages, nil
}

func GetChatMessagesBetweenUsers(db *gorm.DB, requestID, user1ID, user2ID uuid.UUID) ([]Feedback, error) {
	var chatMessages []Feedback
	if err := db.Where("request_id = ? AND (sender_id = ? OR sender_id = ?)", requestID, user1ID, user2ID).Order("created_at asc").Find(&chatMessages).Error; err != nil {
		return nil, err
	}

	return chatMessages, nil
}
