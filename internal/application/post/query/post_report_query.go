package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"time"
)

type GetOnePostReportQuery struct {
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
}

type GetManyPostReportQuery struct {
	Reason       string
	CreatedAt    time.Time
	SortBy       string
	IsDescending bool
	Limit        int
	Page         int
}

type PostReportQueryResult struct {
	PostReport     *common.PostReportResult
	ResultCode     int
	HttpStatusCode int
}

type PostReportQueryListResult struct {
	PostReports    []*common.PostReportShortVerResult
	ResultCode     int
	HttpStatusCode int
	PagingResponse *response.PagingResponse
}
