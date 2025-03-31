package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID                 uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name               string                `gorm:"type:varchar(30);"`
	Image              string                `gorm:"type:varchar(255);default:null"`
	CreatedAt          time.Time             `gorm:"autoCreateTime"`
	UpdatedAt          time.Time             `gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt        `gorm:"index"`
	ConversationDetail []*ConversationDetail `gorm:"foreignKey:ConversationId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
