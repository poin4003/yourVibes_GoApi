package model

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID              uuid.UUID           `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserId          uuid.UUID           `json:"user_id" gorm:"type:uuid;not null"`
	User            User                `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentId        *uuid.UUID          `json:"parent_id" gorm:"type:uuid;default:null"`
	ParentPost      *Post               `json:"parent_post" gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content         string              `json:"content" gorm:"type:text;not null"`
	LikeCount       int                 `json:"like_count" gorm:"type:int;default:0"`
	CommentCount    int                 `json:"comment_count" gorm:"type:int;default:0"`
	Privacy         consts.PrivacyLevel `json:"privacy" gorm:"type:varchar(20);default:'public'"`
	Location        string              `json:"location" gorm:"type:varchar(255);default:null"`
	IsAdvertisement bool                `json:"is_advertisement" gorm:"type:boolean;default:false"`
	Status          bool                `json:"status" gorm:"default:true"`
	CreatedAt       time.Time           `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time           `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt      `json:"deleted_at" gorm:"index"`
	Media           []Media             `json:"media" gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment         []Comment           `json:"comment" gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
