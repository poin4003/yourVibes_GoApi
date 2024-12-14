package common

import (
	"github.com/google/uuid"
	"time"
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

type CommentForReportResult struct {
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
