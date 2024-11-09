package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Media struct {
	ID        uint           `gorm:"type:int;auto_increment;primary_key"`
	PostId    uuid.UUID      `gorm:"type:uuid;not null"`
	MediaUrl  string         `gorm:"type:varchar(255);not null"`
	Status    bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
