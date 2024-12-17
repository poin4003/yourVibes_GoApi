package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"time"
)

type GetOneCommentReportQuery struct {
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
}

type GetManyCommentReportQuery struct {
	Reason       string
	CreatedAt    time.Time
	Status       *bool
	SortBy       string
	IsDescending bool
	Limit        int
	Page         int
}

type CommentReportQueryResult struct {
	CommentReport  *common.CommentReportResult
	ResultCode     int
	HttpStatusCode int
}

type CommentReportQueryListResult struct {
	CommentReports []*common.CommentReportShortVerResult
	ResultCode     int
	HttpStatusCode int
	PagingResponse *response.PagingResponse
}
