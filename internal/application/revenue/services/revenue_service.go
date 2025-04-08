package services

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/revenue/query"
)

type (
	IRevenue interface {
		GetMonthlyRevenue(ctx context.Context, query *query.GetMonthlyRevenueQuery) (result *query.GetMonthlyRevenueQueryResult, err error)
		GetSystemStats(ctx context.Context, query *query.GetSystemStatsQuery) (result *query.GetSystemStatsQueryResult, err error)
	}
)
