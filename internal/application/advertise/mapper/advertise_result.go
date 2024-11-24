package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	advertise_validator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/validator"
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

	return &common.AdvertiseWithBillResult{
		ID:        advertise.ID,
		PostId:    advertise.PostId,
		StartDate: advertise.StartDate,
		EndDate:   advertise.EndDate,
		CreatedAt: advertise.CreatedAt,
		UpdatedAt: advertise.UpdatedAt,
		Bill:      NewBillWithoutAdvertiseResultFromEntity(advertise.Bill),
	}
}

func NewAdvertiseWithoutBillResultFromEntity(
	advertise *entities.Advertise,
) *common.AdvertiseWithoutBillResult {
	if advertise == nil {
		return nil
	}

	return &common.AdvertiseWithoutBillResult{
		ID:        advertise.ID,
		PostId:    advertise.PostId,
		StartDate: advertise.StartDate,
		EndDate:   advertise.EndDate,
		CreatedAt: advertise.CreatedAt,
		UpdatedAt: advertise.UpdatedAt,
	}
}
