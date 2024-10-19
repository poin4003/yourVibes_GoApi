package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	PostId        uuid.UUID      `json:"post_id" gorm:"type:uuid;not null"`
	UserId        uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	User          User           `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentId      *uuid.UUID     `json:"parent_id" gorm:"type:uuid;default:null"`
	ParentComment *Comment       `json:"parent_comment" gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content       string         `json:"content" gorm:"type:text;not null"`
	LikeCount     int            `json:"like_count" gorm:"type:int;default:0"`
	CommentLeft   int            `json:"comment_left" gorm:"type:int;default:0"`
	CommentRight  int            `json:"comment_right" gorm:"type:int;default:0"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
