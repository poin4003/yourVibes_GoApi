package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Advertise struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	PostId    uuid.UUID      `gorm:"type:uuid;not null;index:idx_advertises_post_id"`
	Post      Post           `gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StartDate time.Time      `gorm:"not null;index:idx_advertises_start_date"`
	EndDate   time.Time      `gorm:"not null;index:idx_advertises_end_date"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Bill      Bill           `gorm:"foreignKey:AdvertiseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
