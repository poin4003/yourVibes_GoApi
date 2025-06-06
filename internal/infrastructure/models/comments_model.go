package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	PostId          uuid.UUID      `gorm:"type:uuid;not null;index:idx_comments_post_id_status,priority:1"`
	Post            *Post          `gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId          uuid.UUID      `gorm:"type:uuid;not null"`
	User            User           `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentId        *uuid.UUID     `gorm:"type:uuid;default:null;index:idx_comments_parent_id"`
	ParentComment   *Comment       `gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content         string         `gorm:"type:text;not null"`
	LikeCount       int            `gorm:"type:int;default:0"`
	RepCommentCount int            `gorm:"type:int;default:0"`
	CreatedAt       time.Time      `gorm:"autoCreateTime;index:idx_comments_created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Status          bool           `gorm:"default:true;index:idx_comments_post_id_status,priority:2"`
}
