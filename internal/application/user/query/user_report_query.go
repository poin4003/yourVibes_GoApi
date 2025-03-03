package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetOneUserReportQuery struct {
	UserId         uuid.UUID
	ReportedUserId uuid.UUID
}

type GetManyUserReportQuery struct {
	Reason            string
	UserEmail         string
	ReportedUserEmail string
	AdminEmail        string
	FromDate          time.Time
	ToDate            time.Time
	CreatedAt         time.Time
	Status            *bool
	SortBy            string
	IsDescending      bool
	Limit             int
	Page              int
}

type UserReportQueryResult struct {
	UserReport *common.UserReportResult
}

type UserReportQueryListResult struct {
	UserReports    []*common.UserReportShortVerResult
	PagingResponse *response.PagingResponse
}
