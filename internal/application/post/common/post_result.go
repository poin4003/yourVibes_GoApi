package common

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type PostResultWithLiked struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *UserResult
	ParentId        *uuid.UUID
	ParentPost      *PostResult
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
	IsLiked         bool
}

type PostResult struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *UserResult
	ParentId        *uuid.UUID
	ParentPost      *PostResultWithLiked
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
