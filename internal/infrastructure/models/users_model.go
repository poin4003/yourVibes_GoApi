package models

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FamilyName   string              `gorm:"type:varchar(255);not null"`
	Name         string              `gorm:"type:varchar(255);not null"`
	Email        string              `gorm:"type:varchar(50);unique;not null"`
	Password     *string             `gorm:"type:varchar(255);"`
	PhoneNumber  *string             `gorm:"type:varchar(15);"`
	Birthday     *time.Time          `gorm:"type:timestamptz;"`
	AvatarUrl    string              `gorm:"type:varchar(255);default:null"`
	CapwallUrl   string              `gorm:"type:varchar(255);default:null"`
	Privacy      consts.PrivacyLevel `gorm:"type:varchar(20);default:'public'"`
	Biography    string              `gorm:"type:text;default:null"`
	AuthType     consts.AuthType     `gorm:"type:varchar(10);default:'local'"`
	AuthGoogleId *string             `gorm:"type:varchar(255);default:null"`
	PostCount    int                 `gorm:"type:int;default:0"`
	FriendCount  int                 `gorm:"type:int;default:0"`
	Status       bool                `gorm:"default:true"`
	CreatedAt    time.Time           `gorm:"autoCreateTime"`
	UpdatedAt    time.Time           `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt      `gorm:"index"`
	Setting      Setting             `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
