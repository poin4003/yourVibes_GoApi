package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type PostForReport struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *UserForReport
	ParentId        *uuid.UUID
	ParentPost      *PostForReport
	Content         string
	LikeCount       int
	CommentCount    int
	Privacy         consts.PrivacyLevel
	Location        string
	IsAdvertisement consts.AdvertiseStatus
	Status          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Media           []*Media
}
