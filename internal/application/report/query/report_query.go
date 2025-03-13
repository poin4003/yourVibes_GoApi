package query

import (
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type GetOneReportQuery struct {
	ReportType consts.ReportType
	ReportedId uuid.UUID
}

type GetManyReportQuery struct {
	ReportType        consts.ReportType
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

type ReportQueryResult struct {
	Type          consts.ReportType
	UserReport    *common.UserReportResult
	PostReport    *common.PostReportResult
	CommentReport *common.CommentReportResult
}

type ReportQueryListResult struct {
	Type           consts.ReportType
	UserReports    []*common.UserReportShortVerResult
	PostReports    []*common.PostReportShortVerResult
	CommentReports []*common.CommentReportShortVerResult
	PagingResponse *response.PagingResponse
}
