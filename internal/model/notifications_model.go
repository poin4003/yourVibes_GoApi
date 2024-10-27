package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	ID               uint           `json:"id" gorm:"type:int;auto_increment;primary_key"`
	From             string         `json:"from" gorm:"type:varchar(50);"`
	FromUrl          string         `json:"from_url" gorm:"type:varchar(255);"`
	UserId           uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	User             User           `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	NotificationType string         `json:"notification_type" gorm:"type:varchar(50);"`
	Status           bool           `json:"status" gorm:"default:true"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
