package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
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

type AdvertiseDetail struct {
	ID     uuid.UUID `json:"id"`
	PostId uuid.UUID `json:"post_id"`
	// Post
	StartDate    time.Time                `json:"start_date"`
	EndDate      time.Time                `json:"end_date"`
	DayRemaining int                      `json:"day_remaining"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	Bill         *BillWithoutAdvertiseDto `json:"bill"`
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
