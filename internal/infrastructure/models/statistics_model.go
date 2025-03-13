package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Statistics struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;not null;"`
	AvertiseId      uuid.UUID      `gorm:"type:uuid;not null"`
	Avertise        Advertise      `gorm:"foreignKey:AdvertiseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reach           int            `gorm:"type:int;default:0"`
	Clicks          int            `gorm:"type:int;default:0"`
	Impression      int            `gorm:"type:int;default:0"`
	AggregationDate time.Time      `gorm:"type:timestamptz;not null"`
	Status          bool           `gorm:"default:true"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
