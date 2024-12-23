package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type PostForAdvertiseDto struct {
	ID     uuid.UUID
	UserId uuid.UUID
	// User            *UserForAdvertiseDto
	ParentId        *uuid.UUID
	ParentPost      *PostForAdvertiseDto
	Content         string
	LikeCount       int
	CommentCount    int
	Privacy         consts.PrivacyLevel
	Location        string
	IsAdvertisement bool
	Status          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	// Media           []*MediaResult
}
