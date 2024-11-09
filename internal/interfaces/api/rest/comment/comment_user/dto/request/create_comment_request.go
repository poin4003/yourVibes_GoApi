package request

import (
	"github.com/google/uuid"
)

type CreateCommentInput struct {
	PostId   uuid.UUID  `json:"post_id" binding:"required"`
	ParentId *uuid.UUID `json:"parent_id,omitempty"`
	Content  string     `json:"content" binding:"required"`
}
