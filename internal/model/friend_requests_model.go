package model

import "github.com/google/uuid"

type FriendRequest struct {
	UserId   uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key;not null"`
	FriendId uuid.UUID `json:"friend_id" gorm:"type:uuid;primary_key;not null"`
	User     User      `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Friend   User      `json:"friend" gorm:"foreignKey:FriendId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
