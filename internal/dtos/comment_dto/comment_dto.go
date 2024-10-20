package comment_dto

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/user_dto"
	"time"
)

type CommentDto struct {
	ID        uuid.UUID                `json:"id"`
	PostId    uuid.UUID                `json:"post_id"`
	UserId    uuid.UUID                `json:"user_id"`
	ParentId  *uuid.UUID               `json:"parent_id"`
	Content   string                   `json:"content"`
	LikeCount int                      `json:"like_count"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
	User      user_dto.UserDtoShortVer `json:"user"`
}
