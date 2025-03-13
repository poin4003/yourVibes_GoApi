package common

import (
	"time"

	"github.com/google/uuid"
)

type CommentReportResult struct {
	ReportId          uuid.UUID
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
	ReportId          uuid.UUID
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
	AdminId           *uuid.UUID
	Reason            string
	UserEmail         string
	AdminEmail        *string
	Status            bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
