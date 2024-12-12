package models

import (
	"github.com/google/uuid"
	"time"
)

type CommentReport struct {
	UserId            uuid.UUID `gorm:"type:uuid;primary_key;not null;"`
	ReportedCommentId uuid.UUID `gorm:"type:uuid;primary_key;not null;"`
	AdminId           uuid.UUID `gorm:"type:uuid;not null;"`
	User              User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReportedComment   Comment   `gorm:"foreignKey:ReportedCommentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Admin             Admin     `gorm:"foreignKey:AdminId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reason            string    `gorm:"type:varchar(255);not null"`
	Status            bool      `gorm:"default:true"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
	DeletedAt         time.Time `gorm:"index"`
}
