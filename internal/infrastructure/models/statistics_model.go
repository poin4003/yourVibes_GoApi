package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Statistics struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;not null;"`
	PostId     uuid.UUID      `gorm:"type:uuid;not null;index:idx_statistics_post_id_created_at,priority:1"`
	Reach      int            `gorm:"type:int;default:0"`
	Clicks     int            `gorm:"type:int;default:0"`
	Impression int            `gorm:"type:int;default:0"`
	Status     bool           `gorm:"default:true"`
	CreatedAt  time.Time      `gorm:"autoCreateTime;index:idx_statistics_post_id_created_at,priority:2"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
