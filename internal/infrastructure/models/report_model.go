package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"gorm.io/gorm"
)

type Report struct {
	ID        uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserId    uuid.UUID         `gorm:"type:uuid;not null"`
	AdminId   *uuid.UUID        `gorm:"type:uuid;"`
	User      *User             `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Admin     *Admin            `gorm:"foreignKey:AdminId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reason    string            `gorm:"type:varchar(255);not null"`
	Type      consts.ReportType `gorm:"type:smallint;default:0"`
	Status    bool              `gorm:"default:false"`
	CreatedAt time.Time         `gorm:"autoCreateTime"`
	UpdatedAt time.Time         `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt    `gorm:"index"`
}

type UserReport struct {
	ReportID       uuid.UUID `gorm:"type:uuid;primary_key;not null;"`
	Report         *Report   `gorm:"foreignKey:ReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReportedUserId uuid.UUID `gorm:"type:uuid;primary_key;not null;"`
	ReportedUser   *User     `gorm:"foreignKey:ReportedUserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type PostReport struct {
	ReportID       uuid.UUID `gorm:"type:uuid;primary_key;not null;"`
	Report         *Report   `gorm:"foreignKey:ReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReportedPostId uuid.UUID `gorm:"type:uuid;primary_key;not null;"`
	ReportedPost   *Post     `gorm:"foreignKey:ReportedPostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CommentReport struct {
	ReportID          uuid.UUID `gorm:"type:uuid;primary_key;not null;"`
	Report            *Report   `gorm:"foreignKey:ReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReportedCommentId uuid.UUID `gorm:"type:uuid;primary_key;not null;"`
	ReportedComment   *Comment  `gorm:"foreignKey:ReportedCommentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
