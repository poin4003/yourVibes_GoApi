package entities

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type PostForAdvertise struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *UserForAdvertise
	ParentId        *uuid.UUID
	ParentPost      *PostForAdvertise
	Content         string
	LikeCount       int
	CommentCount    int
	Privacy         consts.PrivacyLevel
	Location        string
	IsAdvertisement bool
	Status          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Media           []*Media
}