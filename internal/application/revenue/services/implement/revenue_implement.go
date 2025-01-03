package implement

import (
	"context"
	revenueQuery "github.com/poin4003/yourVibes_GoApi/internal/application/revenue/query"
	advertiseRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type sRevenue struct {
	billRepo advertiseRepo.IBillRepository
	userRepo advertiseRepo.IUserRepository
	postRepo advertiseRepo.IPostRepository
}

func NewRevenueImplement(
	billRepo advertiseRepo.IBillRepository,
	userRepo advertiseRepo.IUserRepository,
	postRepo advertiseRepo.IPostRepository,
) *sRevenue {
	return &sRevenue{
		billRepo: billRepo,
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (s *sRevenue) GetMonthlyRevenue(
	ctx context.Context,
	query *revenueQuery.GetMonthlyRevenueQuery,
) (result *revenueQuery.GetMonthlyRevenueQueryResult, err error) {
	result = &revenueQuery.GetMonthlyRevenueQueryResult{}
	result.RevenueList = nil
	result.MonthList = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Get list of monthly revenue
	monthList, revenueList, err := s.billRepo.GetMonthlyRevenue(ctx, query.Date)
	if err != nil {
		return result, err
	}

	result.MonthList = monthList
	result.RevenueList = revenueList
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sRevenue) GetSystemStats(
	ctx context.Context,
	query *revenueQuery.GetSystemStatsQuery,
) (result *revenueQuery.GetSystemStatsQueryResult, err error) {
	result = &revenueQuery.GetSystemStatsQueryResult{}
	result.PreviousMonthsRevenue = 0
	result.PreviousDaysRevenue = 0
	result.TotalCountOfUsers = 0
	result.TotalCountOfPosts = 0
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Get previous month revenue
	previousMonthDate := query.Date.AddDate(0, -1, 0)
	previousMonthsRevenue, err := s.billRepo.GetRevenueForMonth(ctx, previousMonthDate)
	if err != nil {
		return result, err
	}

	// 2. Get previous day revenue
	previousDayDate := query.Date.AddDate(0, 0, -1)
	previousDaysRevenue, err := s.billRepo.GetRevenueForDay(ctx, previousDayDate)
	if err != nil {
		return result, err
	}

	// 3. Get total count of user
	totalCountOfUser, err := s.userRepo.GetTotalUserCount(ctx)
	if err != nil {
		return result, err
	}

	// 4. Get total count of post
	totalCountOfPost, err := s.postRepo.GetTotalPostCount(ctx)
	if err != nil {
		return result, err
	}

	result.PreviousMonthsRevenue = previousMonthsRevenue
	result.PreviousDaysRevenue = previousDaysRevenue
	result.TotalCountOfUsers = totalCountOfUser
	result.TotalCountOfPosts = totalCountOfPost
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
