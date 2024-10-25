package post_dto

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/user_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"gorm.io/gorm"
	"time"
)

type PostDto struct {
	ID              uuid.UUID                `json:"id"`
	ParentId        *uuid.UUID               `json:"parent_id"`
	ParentPost      *ParentPostDto           `json:"parent_post"`
	Content         string                   `json:"content"`
	LikeCount       int                      `json:"like_count"`
	CommentCount    int                      `json:"comment_count"`
	Privacy         consts.PrivacyLevel      `json:"privacy"`
	Location        string                   `json:"location"`
	IsAdvertisement bool                     `json:"is_advertisement"`
	Status          bool                     `json:"status"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	DeletedAt       gorm.DeletedAt           `json:"deleted_at"`
	IsLiked         bool                     `json:"is_liked"`
	UserId          uuid.UUID                `json:"user_id"`
	User            user_dto.UserDtoShortVer `json:"user"`
	Media           []model.Media            `json:"media"`
}

type ParentPostDto struct {
	ID              uuid.UUID                `json:"id"`
	Content         string                   `json:"content"`
	LikeCount       int                      `json:"like_count"`
	CommentCount    int                      `json:"comment_count"`
	Privacy         consts.PrivacyLevel      `json:"privacy"`
	Location        string                   `json:"location"`
	IsAdvertisement bool                     `json:"is_advertisement"`
	Status          bool                     `json:"status"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	DeletedAt       gorm.DeletedAt           `json:"deleted_at"`
	IsLiked         bool                     `json:"is_liked"`
	UserId          uuid.UUID                `json:"user_id"`
	User            user_dto.UserDtoShortVer `json:"user"`
	Media           []model.Media            `json:"media"`
}

type UpdatedPostDto struct {
	ID              uuid.UUID                `json:"id"`
	ParentId        *uuid.UUID               `json:"parent_id"`
	ParentPost      *model.Post              `json:"parent_post"`
	Content         string                   `json:"content"`
	LikeCount       int                      `json:"like_count"`
	CommentCount    int                      `json:"comment_count"`
	Privacy         consts.PrivacyLevel      `json:"privacy"`
	Location        string                   `json:"location"`
	IsAdvertisement bool                     `json:"is_advertisement"`
	Status          bool                     `json:"status"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	DeletedAt       gorm.DeletedAt           `json:"deleted_at"`
	UserId          uuid.UUID                `json:"user_id"`
	User            user_dto.UserDtoShortVer `json:"user"`
	Media           []model.Media            `json:"media"`
}

type NewPostDto struct {
	ID              uuid.UUID           `json:"id"`
	ParentId        *uuid.UUID          `json:"parent_id"`
	ParentPost      *model.Post         `json:"parent_post"`
	Content         string              `json:"content"`
	LikeCount       int                 `json:"like_count"`
	CommentCount    int                 `json:"comment_count"`
	Privacy         consts.PrivacyLevel `json:"privacy"`
	Location        string              `json:"location"`
	IsAdvertisement bool                `json:"is_advertisement"`
	Status          bool                `json:"status"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	DeletedAt       gorm.DeletedAt      `json:"deleted_at"`
	UserId          uuid.UUID           `json:"user_id"`
}
