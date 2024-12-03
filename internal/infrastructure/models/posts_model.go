package models

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID              uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserId          uuid.UUID           `gorm:"type:uuid;not null"`
	User            User                `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentId        *uuid.UUID          `gorm:"type:uuid;default:null"`
	ParentPost      *Post               `gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content         string              `gorm:"type:text;not null"`
	LikeCount       int                 `gorm:"type:int;default:0"`
	CommentCount    int                 `gorm:"type:int;default:0"`
	Privacy         consts.PrivacyLevel `gorm:"type:varchar(20);default:'public'"`
	Location        string              `gorm:"type:varchar(255);default:null"`
	IsAdvertisement bool                `gorm:"type:boolean;default:false"`
	Status          bool                `gorm:"default:true"`
	CreatedAt       time.Time           `gorm:"autoCreateTime"`
	UpdatedAt       time.Time           `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt      `gorm:"index"`
	Media           []Media             `gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type PostWithLiked struct {
	ID              uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserId          uuid.UUID           `gorm:"type:uuid;not null"`
	User            User                `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentId        *uuid.UUID          `gorm:"type:uuid;default:null"`
	ParentPost      *Post               `gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content         string              `gorm:"type:text;not null"`
	LikeCount       int                 `gorm:"type:int;default:0"`
	CommentCount    int                 `gorm:"type:int;default:0"`
	Privacy         consts.PrivacyLevel `gorm:"type:varchar(20);default:'public'"`
	Location        string              `gorm:"type:varchar(255);default:null"`
	IsAdvertisement bool                `gorm:"type:boolean;default:false"`
	Status          bool                `gorm:"default:true"`
	CreatedAt       time.Time           `gorm:"autoCreateTime"`
	UpdatedAt       time.Time           `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt      `gorm:"index"`
	Media           []Media             `gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsLiked         bool
}
