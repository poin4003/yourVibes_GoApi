package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"time"
)

type GetOneUserReportQuery struct {
	UserId         uuid.UUID
	ReportedUserId uuid.UUID
}

type GetManyUserReportQuery struct {
	Reason       string
	CreatedAt    time.Time
	SortBy       string
	IsDescending bool
	Limit        int
	Page         int
}

type UserReportQueryResult struct {
	UserReport     *common.UserReportResult
	ResultCode     int
	HttpStatusCode int
}

type UserReportQueryListResult struct {
	UserReports    []*common.UserReportResult
	ResultCode     int
	HttpStatusCode int
	PagingResponse *response.PagingResponse
}