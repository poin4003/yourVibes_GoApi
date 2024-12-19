package common

import (
	"github.com/google/uuid"
	"time"
)

type CommentReportResult struct {
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
	AdminId           *uuid.UUID
	User              *UserForReportResult
	ReportedComment   *CommentForReportResult
	Admin             *AdminResult
	Post              *PostForReportResult
	Reason            string
	Status            bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CommentReportShortVerResult struct {
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
	AdminId           *uuid.UUID
	Reason            string
	Status            bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
