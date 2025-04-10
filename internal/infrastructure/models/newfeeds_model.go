package models

import "github.com/google/uuid"

type NewFeed struct {
	UserId uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	PostId uuid.UUID `gorm:"type:uuid;primary_key;not null;index:idx_new_feeds_post_id"`
	User   User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post   Post      `gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
