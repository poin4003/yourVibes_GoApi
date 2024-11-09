package models

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	ID               uint                    `gorm:"type:int;auto_increment;primary_key"`
	From             string                  `gorm:"type:varchar(50);"`
	FromUrl          string                  `gorm:"type:varchar(255);"`
	UserId           uuid.UUID               `gorm:"type:uuid;not null"`
	User             User                    `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	NotificationType consts.NotificationType `gorm:"type:varchar(50);"`
	ContentId        string                  `gorm:"type:varchar(50);"`
	Content          string                  `gorm:"type:text;not null"`
	Status           bool                    `gorm:"default:true"`
	CreatedAt        time.Time               `gorm:"autoCreateTime"`
	UpdatedAt        time.Time               `gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt          `gorm:"index"`
}
