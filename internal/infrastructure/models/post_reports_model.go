package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PostReport struct {
	UserId         uuid.UUID      `gorm:"type:uuid;primary_key;not null;"`
	ReportedPostId uuid.UUID      `gorm:"type:uuid;primary_key;not null;"`
	AdminId        *uuid.UUID     `gorm:"type:uuid;"`
	User           User           `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReportedPost   Post           `gorm:"foreignKey:ReportedPostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Admin          Admin          `gorm:"foreignKey:AdminId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reason         string         `gorm:"type:varchar(255);not null"`
	Status         bool           `gorm:"default:false"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
