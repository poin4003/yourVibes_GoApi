package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Bill struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	AdvertiseId uuid.UUID      `gorm:"type:uuid;not null"`
	Advertise   *Advertise     `gorm:"foreignKey:AdvertiseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Price       float64        `gorm:"type:decimal(10,2);default:0.0"`
	Vat         float64        `gorm:"type:decimal(10,2);default:0.0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Status      bool           `gorm:"default:true"`
}
