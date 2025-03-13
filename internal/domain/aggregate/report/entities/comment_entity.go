package entities

import (
	"time"

	"github.com/google/uuid"
)

type CommentForReport struct {
	ID              uuid.UUID
	PostId          uuid.UUID
	UserId          uuid.UUID
	User            *UserForReport
	ParentId        *uuid.UUID
	Content         string
	LikeCount       int
	RepCommentCount int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Status          bool
}
