package services

import (
	"context"

	"github.com/poin4003/yourVibes_GoApi/internal/application/report/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/query"
)

type (
	IReport interface {
		CreateReport(ctx context.Context, command *command.CreateReportCommand) error
		HandleReport(ctx context.Context, command *command.HandleReportCommand) (err error)
		DeleteReport(ctx context.Context, command *command.DeleteReportCommand) (err error)
		Activate(ctx context.Context, command *command.ActivateCommand) (err error)
		GetDetailReport(ctx context.Context, query *query.GetOneReportQuery) (result *query.ReportQueryResult, err error)
		GetManyReport(ctx context.Context, query *query.GetManyReportQuery) (result *query.ReportQueryListResult, err error)
	}
)

var (
	localReport IReport
)

func Report() IReport {
	if localReport == nil {
		panic("repository_implement localReport not found for interface IReport")
	}

	return localReport
}

func InitReport(i IReport) {
	localReport = i
}
