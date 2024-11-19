package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Bill struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	AdvertiesId uuid.UUID      `gorm:"type:uuid;not null"`
	Price       float64        `gorm:"type:decimal(10,2);default:0.0"`
	Vat         float64        `gorm:"type:decimal(10,2);default:0.0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Status      bool           `gorm:"default:true"`
}
