package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Media struct {
	ID          uint           `json:"id" gorm:"type:int;auto_increment;primary_key"`
	PostId      uuid.UUID      `json:"post_id" gorm:"type:uuid;not null"`
	Post        Post           `json:"post" gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MediaUrl    string         `json:"media_url" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:text;not null"`
	Status      bool           `json:"status" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
