package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetOneCommentReportQuery struct {
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
}

type GetManyCommentReportQuery struct {
	Reason       string
	UserEmail    string
	AdminEmail   string
	FromDate     time.Time
	ToDate       time.Time
	CreatedAt    time.Time
	Status       *bool
	SortBy       string
	IsDescending bool
	Limit        int
	Page         int
}

type CommentReportQueryResult struct {
	CommentReport *common.CommentReportResult
}

type CommentReportQueryListResult struct {
	CommentReports []*common.CommentReportShortVerResult
	PagingResponse *response.PagingResponse
}
