package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToBillModel(bill *entities.Bill) *models.Bill {
	b := &models.Bill{
		AdvertiseId: bill.AdvertiseId,
		Price:       bill.Price,
		CreatedAt:   bill.CreatedAt,
		UpdatedAt:   bill.UpdateAt,
		Status:      bill.Status,
		VoucherId:   bill.VoucherId,
	}
	b.ID = bill.ID

	return b
}

func FromBillModel(b *models.Bill) *entities.Bill {
	var advertise = &entities.Advertise{
		ID:        b.Advertise.ID,
		PostId:    b.Advertise.PostId,
		StartDate: b.Advertise.StartDate,
		EndDate:   b.Advertise.EndDate,
		CreatedAt: b.Advertise.CreatedAt,
		UpdatedAt: b.Advertise.UpdatedAt,
	}

	var bill = &entities.Bill{
		AdvertiseId: b.AdvertiseId,
		Advertise:   advertise,
		Price:       b.Price,
		CreatedAt:   b.CreatedAt,
		UpdateAt:    b.UpdatedAt,
		Status:      b.Status,
	}
	bill.ID = b.ID

	return bill
}
