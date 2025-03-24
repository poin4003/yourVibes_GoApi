package common

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type PostForAdvertiseResult struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *UserForAdvertiseResult
	ParentId        *uuid.UUID
	ParentPost      *PostForAdvertiseResult
	Content         string
	LikeCount       int
	CommentCount    int
	Privacy         consts.PrivacyLevel
	Location        string
	IsAdvertisement consts.AdvertiseStatus
	Status          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Media           []*MediaResult
}
