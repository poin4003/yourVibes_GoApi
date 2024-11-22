package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Advertise struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	PostId    uuid.UUID      `gorm:"type:uuid;not null"`
	Post      Post           `gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StartDate time.Time      `gorm:"not null"`
	EndDate   time.Time      `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Bill      Bill           `gorm:"foreignKey:AdvertiseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
