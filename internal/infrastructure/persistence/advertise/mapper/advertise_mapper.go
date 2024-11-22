package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToAdvertiseModel(advertise *entities.Advertise) *models.Advertise {
	a := &models.Advertise{
		PostId:    advertise.PostId,
		StartDate: advertise.StartDate,
		EndDate:   advertise.EndDate,
		CreatedAt: advertise.CreatedAt,
		UpdatedAt: advertise.UpdatedAt,
	}
	a.ID = advertise.ID

	return a
}

func FromAdvertiseModel(a *models.Advertise) *entities.Advertise {
	var advertise = &entities.Advertise{
		PostId:    a.PostId,
		StartDate: a.StartDate,
		EndDate:   a.EndDate,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
	advertise.ID = a.ID

	return advertise
}
