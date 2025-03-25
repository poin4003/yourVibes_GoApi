package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"time"
)

type AdvertiseWithBillDto struct {
	ID           uuid.UUID                `json:"id"`
	PostId       uuid.UUID                `json:"post_id"`
	StartDate    time.Time                `json:"start_date"`
	EndDate      time.Time                `json:"end_date"`
	DayRemaining int                      `json:"day_remaining"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	Bill         *BillWithoutAdvertiseDto `json:"bill"`
}

type AdvertiseWithStatisticDto struct {
	ID              uuid.UUID                `json:"id"`
	PostId          uuid.UUID                `json:"post_id"`
	StartDate       time.Time                `json:"start_date"`
	EndDate         time.Time                `json:"end_date"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	DayRemaining    int                      `json:"day_remaining"`
	Bill            *BillWithoutAdvertiseDto `json:"bill"`
	Post            *PostForAdvertiseDto     `json:"post"`
	TotalReach      int                      `json:"total_reach"`
	TotalClicks     int                      `json:"total_clicks"`
	TotalImpression int                      `json:"total_impression"`
	Statistics      []*StatisticDto          `json:"statistics"`
}

func ToAdvertiseWithBillDto(
	advertiseResult common.AdvertiseWithBillResult,
) *AdvertiseWithBillDto {
	advertiseDto := &AdvertiseWithBillDto{
		ID:           advertiseResult.ID,
		PostId:       advertiseResult.PostId,
		StartDate:    advertiseResult.StartDate,
		EndDate:      advertiseResult.EndDate,
		CreatedAt:    advertiseResult.CreatedAt,
		DayRemaining: advertiseResult.DayRemaining,
		UpdatedAt:    advertiseResult.UpdatedAt,
		Bill:         ToBillWithoutAdvertiseDto(*advertiseResult.Bill),
	}

	return advertiseDto
}

func ToAdvertiseWithStatisticDto(
	advertiseResult common.AdvertiseForStatisticResult,
) *AdvertiseWithStatisticDto {
	if advertiseResult.Statistics == nil {
		return nil
	}

	var statistics []*StatisticDto
	for _, stat := range advertiseResult.Statistics {
		statistics = append(statistics, ToStatisticDto(stat))
	}

	advertiseDto := &AdvertiseWithStatisticDto{
		ID:              advertiseResult.ID,
		PostId:          advertiseResult.PostId,
		StartDate:       advertiseResult.StartDate,
		EndDate:         advertiseResult.EndDate,
		DayRemaining:    advertiseResult.DayRemaining,
		CreatedAt:       advertiseResult.CreatedAt,
		UpdatedAt:       advertiseResult.UpdatedAt,
		Bill:            ToBillWithoutAdvertiseDto(*advertiseResult.Bill),
		Post:            ToPostForAdvertiseDto(advertiseResult.Post),
		TotalReach:      advertiseResult.TotalReach,
		TotalClicks:     advertiseResult.TotalClicks,
		TotalImpression: advertiseResult.TotalImpression,
		Statistics:      statistics,
	}

	return advertiseDto
}
