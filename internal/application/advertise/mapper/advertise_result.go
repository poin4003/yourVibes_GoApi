package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	advertise_validator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/validator"
	"time"
)

func NewAdvertiseWithBillResultFromValidateEntity(
	advertise *advertise_validator.ValidateAdvertise,
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
