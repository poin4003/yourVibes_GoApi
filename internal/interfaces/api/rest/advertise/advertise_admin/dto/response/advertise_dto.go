package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
)

type AdvertiseWithBillDto struct {
	ID           uuid.UUID                `json:"id"`
	PostId       uuid.UUID                `json:"post_id"`
	UserEmail    string                   `json:"user_email"`
	StartDate    time.Time                `json:"start_date"`
	EndDate      time.Time                `json:"end_date"`
	DayRemaining int                      `json:"day_remaining"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	Bill         *BillWithoutAdvertiseDto `json:"bill"`
}

type AdvertiseDetailDto struct {
	ID           uuid.UUID                `json:"id"`
	PostId       uuid.UUID                `json:"post_id"`
	Post         PostForAdvertiseDto      `json:"post"`
	StartDate    time.Time                `json:"start_date"`
	EndDate      time.Time                `json:"end_date"`
	DayRemaining int                      `json:"day_remaining"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	Bill         *BillWithoutAdvertiseDto `json:"bill"`
}

func ToAdvertiseWithBillDto(
	advertiseResult *common.AdvertiseWithBillResult,
) *AdvertiseWithBillDto {
	advertiseDto := &AdvertiseWithBillDto{
		ID:           advertiseResult.ID,
		PostId:       advertiseResult.PostId,
		UserEmail:    advertiseResult.UserEmail,
		StartDate:    advertiseResult.StartDate,
		EndDate:      advertiseResult.EndDate,
		CreatedAt:    advertiseResult.CreatedAt,
		DayRemaining: advertiseResult.DayRemaining,
		UpdatedAt:    advertiseResult.UpdatedAt,
		Bill:         ToBillWithoutAdvertiseDto(*advertiseResult.Bill),
	}

	return advertiseDto
}

func ToAdvertiseDetail(
	advertiseResult *common.AdvertiseDetail,
) *AdvertiseDetailDto {
	advertiseDto := &AdvertiseDetailDto{
		ID:           advertiseResult.ID,
		PostId:       advertiseResult.PostId,
		Post:         *ToPostForAdvertiseDto(advertiseResult.Post),
		StartDate:    advertiseResult.StartDate,
		EndDate:      advertiseResult.EndDate,
		DayRemaining: advertiseResult.DayRemaining,
		CreatedAt:    advertiseResult.CreatedAt,
		UpdatedAt:    advertiseResult.UpdatedAt,
		Bill:         ToBillWithoutAdvertiseDto(*advertiseResult.Bill),
	}

	return advertiseDto
}
