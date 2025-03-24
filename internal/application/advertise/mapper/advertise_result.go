package mapper

import (
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	advertiseValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/validator"
)

func NewAdvertiseWithBillResultFromValidateEntity(
	advertise *advertiseValidator.ValidateAdvertise,
) *common.AdvertiseWithBillResult {
	return NewAdvertiseWithBillResultFromEntity(&advertise.Advertise)
}

func NewAdvertiseWithBillResultFromEntity(
	advertise *entities.Advertise,
) *common.AdvertiseWithBillResult {
	if advertise == nil {
		return nil
	}

	duration := advertise.EndDate.Sub(time.Now())
	dayRemaining := int(duration.Hours() / 24)

	if dayRemaining == 0 && duration.Hours() > 0 {
		dayRemaining = 1
	}

	if dayRemaining < 0 {
		dayRemaining = 0
	}

	return &common.AdvertiseWithBillResult{
		ID:           advertise.ID,
		PostId:       advertise.PostId,
		UserEmail:    advertise.Post.User.Email,
		StartDate:    advertise.StartDate,
		EndDate:      advertise.EndDate,
		DayRemaining: dayRemaining,
		CreatedAt:    advertise.CreatedAt,
		UpdatedAt:    advertise.UpdatedAt,
		Bill:         NewBillWithoutAdvertiseResultFromEntity(advertise.Bill),
	}
}

func NewAdvertiseWithoutBillResultFromEntity(
	advertise *entities.Advertise,
) *common.AdvertiseWithoutBillResult {
	if advertise == nil {
		return nil
	}

	duration := advertise.EndDate.Sub(time.Now())
	dayRemaining := int(duration.Hours() / 24)

	if dayRemaining == 0 && duration.Hours() > 0 {
		dayRemaining = 1
	}

	if dayRemaining < 0 {
		dayRemaining = 0
	}

	return &common.AdvertiseWithoutBillResult{
		ID:           advertise.ID,
		PostId:       advertise.PostId,
		StartDate:    advertise.StartDate,
		EndDate:      advertise.EndDate,
		DayRemaining: dayRemaining,
		CreatedAt:    advertise.CreatedAt,
		UpdatedAt:    advertise.UpdatedAt,
	}
}

func NewAdvertiseDetailFromEntity(
	advertise *entities.Advertise,
) *common.AdvertiseDetailResult {

	if advertise == nil {
		return nil
	}

	duration := advertise.EndDate.Sub(time.Now())
	dayRemaining := int(duration.Hours() / 24)

	if dayRemaining == 0 && duration.Hours() > 0 {
		dayRemaining = 1
	}

	if dayRemaining < 0 {
		dayRemaining = 0
	}

	return &common.AdvertiseDetailResult{
		ID:           advertise.ID,
		PostId:       advertise.PostId,
		Post:         NewPostResult(advertise.Post),
		StartDate:    advertise.StartDate,
		EndDate:      advertise.EndDate,
		DayRemaining: dayRemaining,
		CreatedAt:    advertise.CreatedAt,
		UpdatedAt:    advertise.UpdatedAt,
		Bill:         NewBillWithoutAdvertiseResultFromEntity(advertise.Bill),
	}
}

func NewAdvertiseDetailAndStatisticResult(
	advertise *entities.AdvertiseForStatistic,
) *common.AdvertiseForStatisticResult {

	if advertise == nil {
		return nil
	}

	duration := advertise.EndDate.Sub(time.Now())
	dayRemaining := int(duration.Hours() / 24)

	if dayRemaining == 0 && duration.Hours() > 0 {
		dayRemaining = 1
	}

	if dayRemaining < 0 {
		dayRemaining = 0
	}

	var statisticResults []*common.StatisticResult
	for _, stat := range advertise.Statistics {
		statisticResults = append(statisticResults, NewStatisticResult(stat))
	}

	return &common.AdvertiseForStatisticResult{
		ID:              advertise.ID,
		PostId:          advertise.PostId,
		StartDate:       advertise.StartDate,
		EndDate:         advertise.EndDate,
		DayRemaining:    dayRemaining,
		CreatedAt:       advertise.CreatedAt,
		UpdatedAt:       advertise.UpdatedAt,
		TotalReach:      advertise.TotalReach,
		TotalClicks:     advertise.TotalClicks,
		TotalImpression: advertise.TotalImpression,
		Statistics:      statisticResults,
	}
}
