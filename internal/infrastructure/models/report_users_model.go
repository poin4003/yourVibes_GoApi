package models

import (
	"github.com/google/uuid"
	"time"
)

type ReportUsersModel struct {
	UserId         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	User           User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReportedUserId uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ReportUser     User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AdminId        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Reason         string    `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	DeletedAt      time.Time `gorm:"index"`
}
