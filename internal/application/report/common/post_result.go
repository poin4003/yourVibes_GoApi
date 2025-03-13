package common

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type PostForReportResult struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *UserForReportResult
	ParentId        *uuid.UUID
	ParentPost      *PostForReportResult
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
