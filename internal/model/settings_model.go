package model

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"gorm.io/gorm"
	"time"
)

type Setting struct {
	ID        uint            `json:"id" gorm:"type:int;auto_increment;primary_key"`
	UserId    uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	Language  consts.Language `json:"language" gorm:"type:varchar(10);default:'vi'"`
	Status    bool            `json:"status" gorm:"default:true"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt  `json:"deleted_at" gorm:"index"`
}
