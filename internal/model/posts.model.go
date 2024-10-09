package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserId       string         `json:"family_name" gorm:"type:varchar(255);not null"`
	User         User           `json:""`
	Name         string         `json:"name" gorm:"type:varchar(255);not null"`
	Email        string         `json:"email" gorm:"type:varchar(50);unique;not null"`
	Password     string         `json:"password" gorm:"type:varchar(255);not null"`
	PhoneNumber  string         `json:"phone_number" gorm:"type:varchar(15);not null"`
	Birthday     time.Time      `json:"birthday" gorm:"type:timestamptz;not null"`
	AvatarUrl    string         `json:"avatar_url" gorm:"type:varchar(255);default:null"`
	CapwallUrl   string         `json:"capwall_url" gorm:"type:varchar(255);default:null"`
	Privacy      string         `json:"privacy" gorm:"type:varchar(20);default:'public"`
	AuthType     string         `json:"auth_type" gorm:"type:varchar(10);default:'local'"`
	AuthGoogleId string         `json:"auth_google_id" gorm:"type:varchar(255);default:null"`
	Status       bool           `json:"status" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
