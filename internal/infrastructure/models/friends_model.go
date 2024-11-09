package models

import (
	"github.com/google/uuid"
)

type Friend struct {
	UserId   uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	FriendId uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	User     User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Friend   User      `gorm:"foreignKey:FriendId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
