package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetOnePostReportQuery struct {
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
}

type GetManyPostReportQuery struct {
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

type PostReportQueryResult struct {
	PostReport *common.PostReportResult
}

type PostReportQueryListResult struct {
	PostReports    []*common.PostReportShortVerResult
	PagingResponse *response.PagingResponse
}
