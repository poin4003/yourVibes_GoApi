package common

import (
	"time"

	"github.com/google/uuid"
)

type CommentResult struct {
	ID              uuid.UUID
	PostId          uuid.UUID
	UserId          uuid.UUID
	User            *UserResult
	ParentId        *uuid.UUID
	Content         string
	LikeCount       int
	RepCommentCount int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Status          bool
}

type CommentResultWithLiked struct {
	ID              uuid.UUID
	PostId          uuid.UUID
	UserId          uuid.UUID
	User            *UserResult
	ParentId        *uuid.UUID
	Content         string
	LikeCount       int
	RepCommentCount int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Status          bool
	IsLiked         bool
}
