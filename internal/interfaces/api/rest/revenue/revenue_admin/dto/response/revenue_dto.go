package response

import "github.com/poin4003/yourVibes_GoApi/internal/application/revenue/query"

type MonthlyRevenueDto struct {
	MonthList   []string `json:"month_list"`
	RevenueList []int64  `json:"revenue_list"`
}

type SystemStatsDto struct {
	PreviousMonthsRevenue int64 `json:"previous_months_revenue"`
	PreviousDaysRevenue   int64 `json:"previous_days_revenue"`
	TotalCountOfUsers     int   `json:"total_count_of_users"`
	TotalCountOfPosts     int   `json:"total_count_of_posts"`
}

func ToMonthlyRevenueDto(monthRevenueResult *query.GetMonthlyRevenueQueryResult) *MonthlyRevenueDto {
	return &MonthlyRevenueDto{
		MonthList:   monthRevenueResult.MonthList,
		RevenueList: monthRevenueResult.RevenueList,
	}
}

func ToSystemStatsDto(systemStatsResult *query.GetSystemStatsQueryResult) *SystemStatsDto {
	return &SystemStatsDto{
		PreviousMonthsRevenue: systemStatsResult.PreviousMonthsRevenue,
		PreviousDaysRevenue:   systemStatsResult.PreviousDaysRevenue,
		TotalCountOfUsers:     systemStatsResult.TotalCountOfUsers,
		TotalCountOfPosts:     systemStatsResult.TotalCountOfPosts,
	}
}
