package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserId         uuid.UUID      `gorm:"type:uuid;not null"`
	User           User           `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ConversationId uuid.UUID      `gorm:"type:uuid;not null"`
	Conversation   Conversation   `gorm:"foreignKey:ConversationId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentId       *uuid.UUID     `gorm:"type:uuid;default:null"`
	ParentMessage  *Message       `gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content        *string        `gorm:"type:text;"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
