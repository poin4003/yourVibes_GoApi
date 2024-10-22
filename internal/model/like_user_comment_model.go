package model

import "github.com/google/uuid"

type LikeUserComment struct {
	UserId    uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key;not null"`
	CommentId uuid.UUID `json:"comment_id" gorm:"type:uuid;primary_key;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment   Comment   `json:"comment" gorm:"foreignKey:CommentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
