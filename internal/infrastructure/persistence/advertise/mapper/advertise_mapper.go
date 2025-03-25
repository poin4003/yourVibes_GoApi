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
	var post = &entities.PostForAdvertise{
		User: FromUserModel(&a.Post.User),
	}

	var bill = &entities.Bill{
		ID:          a.Bill.ID,
		AdvertiseId: a.Bill.AdvertiseId,
		Price:       a.Bill.Price,
		CreatedAt:   a.Bill.CreatedAt,
		UpdateAt:    a.Bill.UpdatedAt,
		Status:      a.Bill.Status,
	}

	var advertise = &entities.Advertise{
		PostId:    a.PostId,
		Post:      post,
		StartDate: a.StartDate,
		EndDate:   a.EndDate,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Bill:      bill,
	}
	advertise.ID = a.ID

	return advertise
}

func FromAdvertiseModelForAdvertiseDetail(a *models.Advertise) *entities.Advertise {
	var bill = &entities.Bill{
		ID:          a.Bill.ID,
		AdvertiseId: a.Bill.AdvertiseId,
		Price:       a.Bill.Price,
		CreatedAt:   a.Bill.CreatedAt,
		UpdateAt:    a.Bill.UpdatedAt,
		Status:      a.Bill.Status,
	}

	var advertise = &entities.Advertise{
		PostId:    a.PostId,
		Post:      FromPostModel(&a.Post),
		StartDate: a.StartDate,
		EndDate:   a.EndDate,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Bill:      bill,
	}
	advertise.ID = a.ID

	return advertise
}

func FromAdvertiseModelForDetailAndStatistics(
	a *models.Advertise,
	totalReach int,
	totalClicks int,
	totalImpressions int,
	statEntities []*entities.StatisticEntity,
) *entities.AdvertiseForStatistic {
	if a == nil {
		return nil
	}

	var bill = &entities.Bill{
		ID:          a.Bill.ID,
		AdvertiseId: a.Bill.AdvertiseId,
		Price:       a.Bill.Price,
		CreatedAt:   a.Bill.CreatedAt,
		UpdateAt:    a.Bill.UpdatedAt,
		Status:      a.Bill.Status,
	}

	var advertise = &entities.AdvertiseForStatistic{
		PostId:           a.PostId,
		StartDate:        a.StartDate,
		EndDate:          a.EndDate,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
		Bill:             bill,
		PostForAdvertise: FromPostModel(&a.Post),
		TotalReach:       totalReach,
		TotalClicks:      totalClicks,
		TotalImpression:  totalImpressions,
		Statistics:       statEntities,
	}
	advertise.ID = a.ID

	return advertise
}
