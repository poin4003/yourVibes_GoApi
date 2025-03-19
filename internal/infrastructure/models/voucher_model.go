package models

import (
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Voucher struct {
	ID          uuid.UUID            `gorm:"type:uuid;primary_key;not null;"`
	AdminId     *uuid.UUID           `gorm:"type:uuid;"`
	Admin       *Admin               `gorm:"foreignKey:AdminId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string               `gorm:"varchar(50);not null"`
	Description string               `gorm:"text;default:null"`
	Type        consts.VoucherType   `gorm:"default:false"`
	Value       int                  `gorm:"type:int;default:1"`
	Code        string               `gorm:"type:varchar(30);not null"`
	MaxUses     int                  `gorm:"type:int;default:1"`
	Status      consts.VoucherStatus `gorm:"default:true"`
	CreatedAt   time.Time            `gorm:"autoCreateTime"`
	UpdatedAt   time.Time            `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt       `gorm:"index"`
}
