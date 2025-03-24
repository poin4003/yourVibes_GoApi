package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Statistics struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;not null;"`
	PostId     uuid.UUID      `gorm:"type:uuid;not null"`
	Reach      int            `gorm:"type:int;default:0"`
	Clicks     int            `gorm:"type:int;default:0"`
	Impression int            `gorm:"type:int;default:0"`
	Status     bool           `gorm:"default:true"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
