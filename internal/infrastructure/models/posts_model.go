package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"gorm.io/gorm"
)

type Post struct {
	ID              uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserId          uuid.UUID              `gorm:"type:uuid;not null;index:idx_posts_status_user_id,priority:2"`
	User            User                   `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentId        *uuid.UUID             `gorm:"type:uuid;default:null"`
	ParentPost      *Post                  `gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content         string                 `gorm:"type:text;not null"`
	LikeCount       int                    `gorm:"type:int;default:0"`
	CommentCount    int                    `gorm:"type:int;default:0"`
	Privacy         consts.PrivacyLevel    `gorm:"type:varchar(20);default:'public';index:idx_posts_privacy"`
	Location        string                 `gorm:"type:varchar(255);default:null"`
	IsAdvertisement consts.AdvertiseStatus `gorm:"type:smallint;default:0;index:idx_posts_is_advertisement"`
	IsFeaturedPost  bool                   `gorm:"type:boolean;default:false"`
	Status          bool                   `gorm:"type:boolean;default:false;index:idx_posts_status_user_id,priority:1"`
	CreatedAt       time.Time              `gorm:"autoCreateTime;index:idx_posts_created_at"`
	UpdatedAt       time.Time              `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt         `gorm:"index"`
	Media           []Media                `gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Statistics      []Statistics           `gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
