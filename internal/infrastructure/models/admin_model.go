package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FamilyName  string         `gorm:"type:varchar(255);not null"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Email       string         `gorm:"type:varchar(50);unique;not null"`
	Password    string         `gorm:"type:varchar(255);not null"`
	PhoneNumber string         `gorm:"type:varchar(15);not null"`
	IdentityId  string         `gorm:"type:varchar(15);not null"`
	Birthday    time.Time      `gorm:"type:timestamptz;not null"`
	Status      bool           `gorm:"default:true"`
	Role        bool           `gorm:"default:false"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
