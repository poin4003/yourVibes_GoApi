package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bill struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	AdvertiseId uuid.UUID      `gorm:"type:uuid;not null"`
	Advertise   *Advertise     `gorm:"foreignKey:AdvertiseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Price       int            `gorm:"default:0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Status      bool           `gorm:"default:false"`
	VoucherId   uuid.UUID      `gorm:"type:uuid;not null"`
	Voucher     Voucher        `gorm:"foreignKey:VoucherId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
