package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
)

func NewBillWithoutAdvertiseResultFromEntity(
	bill *entities.Bill,
) *common.BillWithoutAdvertiseResult {
	if bill == nil {
		return nil
	}

	return &common.BillWithoutAdvertiseResult{
		ID:          bill.ID,
		AdvertiseId: bill.AdvertiseId,
		Price:       bill.Price,
		CreatedAt:   bill.CreatedAt,
		UpdatedAt:   bill.UpdateAt,
		Status:      bill.Status,
	}
}

func NewBillWithAdvertiseResultFromEntity(
	bill *entities.Bill,
) *common.BillWithAdvertiseResult {
	if bill == nil {
		return nil
	}

	return &common.BillWithAdvertiseResult{
		ID:          bill.ID,
		AdvertiseId: bill.AdvertiseId,
		Advertise:   NewAdvertiseWithoutBillResultFromEntity(bill.Advertise),
		Price:       bill.Price,
		CreatedAt:   bill.CreatedAt,
		UpdatedAt:   bill.UpdateAt,
		Status:      bill.Status,
	}
}
