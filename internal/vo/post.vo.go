package vo

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type CreatePostInput struct {
	UserId   uuid.UUID           `json:"user_id" binding:"required"`
	ParentId *uuid.UUID          `json:"parent_id,omitempty"`
	Title    string              `json:"title" binding:"required"`
	Content  string              `json:"content" binding:"required"`
	Privacy  consts.PrivacyLevel `json:"privacy" binding:"required,privacy_enum"`
	Location string              `json:"location,omitempty"`
}

type UpdatePostInput struct {
	Title    *string              `json:"title" binding:"required"`
	Content  *string              `json:"content" binding:"required"`
	Privacy  *consts.PrivacyLevel `json:"privacy" binding:"required,privacy_enum"`
	Location *string              `json:"location,omitempty"`
}
