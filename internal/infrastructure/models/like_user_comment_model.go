package models

import (
	"github.com/google/uuid"
)

type LikeUserComment struct {
	UserId    uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	CommentId uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	User      User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment   Comment   `gorm:"foreignKey:CommentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
