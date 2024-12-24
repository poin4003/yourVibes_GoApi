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

var (
	localRevenue IRevenue
)

func Revenue() IRevenue {
	if localRevenue == nil {
		panic("service implement local IRevenue not found for interface")
	}
	return localRevenue
}

func InitRevenue(i IRevenue) {
	localRevenue = i
}
