package models

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"gorm.io/gorm"
	"time"
)

type Setting struct {
	ID        uint            `gorm:"type:int;auto_increment;primary_key"`
	UserId    uuid.UUID       `gorm:"type:uuid;not null"`
	Language  consts.Language `gorm:"type:varchar(10);default:'vi'"`
	Status    bool            `gorm:"default:true"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt  `gorm:"index"`
}
