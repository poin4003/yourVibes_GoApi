package query

import "time"

type GetMonthlyRevenueQuery struct {
	Date time.Time
}

type GetSystemStatsQuery struct {
	Date time.Time
}

type GetMonthlyRevenueQueryResult struct {
	MonthList   []string
	RevenueList []int64
}

type GetSystemStatsQueryResult struct {
	PreviousMonthsRevenue int64
	PreviousDaysRevenue   int64
	TotalCountOfUsers     int
	TotalCountOfPosts     int
}
