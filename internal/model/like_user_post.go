package model

import "github.com/google/uuid"

type LikeUserPost struct {
	UserId uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key;not null"`
	PostId uuid.UUID `json:"post_id" gorm:"type:uuid;primary_key;not null"`
	User   User      `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post   Post      `json:"post" gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
