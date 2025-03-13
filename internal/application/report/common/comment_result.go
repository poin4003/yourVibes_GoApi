package common

import (
	"time"

	"github.com/google/uuid"
)

type CommentForReportResult struct {
	ID              uuid.UUID
	PostId          uuid.UUID
	UserId          uuid.UUID
	User            *UserForReportResult
	ParentId        *uuid.UUID
	Content         string
	LikeCount       int
	RepCommentCount int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Status          bool
}
