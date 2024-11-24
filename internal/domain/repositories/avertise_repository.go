package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
)

type (
	IAdvertiseRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Advertise, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Advertise, error)
		CreateOne(ctx context.Context, entity *entities.Advertise) (*entities.Advertise, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.AdvertiseUpdate) (*entities.Advertise, error)
		DeleteOne(ctx context.Context, id uuid.UUID) error
	}
	IBillRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Bill, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Bill, error)
		CreateOne(ctx context.Context, entity *entities.Bill) (*entities.Bill, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.BillUpdate) (*entities.Bill, error)
		DeleteOne(ctx context.Context, id uuid.UUID) error
	}
)

var (
	localAdvertise IAdvertiseRepository
	localBill      IBillRepository
)

func Advertise() IAdvertiseRepository {
	if localAdvertise == nil {
		panic("repository_implement localAdvertise not found for interface IAdvertiseRepository")
	}

	return localAdvertise
}

func Bill() IBillRepository {
	if localBill == nil {
		panic("repository_implement localBill not found for interface IBillRepository")
	}

	return localBill
}

func InitAdvertiseRepository(i IAdvertiseRepository) {
	localAdvertise = i
}

func InitBillRepository(i IBillRepository) {
	localBill = i
}
