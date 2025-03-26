package models

import "github.com/google/uuid"

type ConversationDetail struct {
	UserId         uuid.UUID    `gorm:"type:uuid;primary_key;not null"`
	ConversationId uuid.UUID    `gorm:"type:uuid;primary_key;not null"`
	User           User         `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Conversation   Conversation `gorm:"foreignKey:ConversationId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LastMessStatus bool         `gorm:"default:true"`
	LastMessId     *uuid.UUID   `gorm:"type:uuid;"`
	LastMess       *Message     `gorm:"foreignKey:LastMessId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
