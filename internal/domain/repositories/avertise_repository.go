package repositories

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
)

type (
	IAdvertiseRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Advertise, error)
		GetOne(ctx context.Context, id uuid.UUID) (*entities.Advertise, error)
		GetDetailAndStatisticOfAdvertise(ctx context.Context, id uuid.UUID) (*entities.AdvertiseForStatistic, error)
		GetMany(ctx context.Context, query *query.GetManyAdvertiseQuery) ([]*entities.Advertise, *response.PagingResponse, error)
		CreateOne(ctx context.Context, entity *entities.Advertise) (*entities.Advertise, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.AdvertiseUpdate) (*entities.Advertise, error)
		DeleteOne(ctx context.Context, id uuid.UUID) error
		DeleteMany(ctx context.Context, condition map[string]interface{}) error
		GetLatestAdsByPostId(ctx context.Context, postId uuid.UUID) (*entities.Advertise, error)
		CheckExists(ctx context.Context, postId uuid.UUID) (bool, error)
	}
	IBillRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Bill, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Bill, error)
		CreateOne(ctx context.Context, entity *entities.Bill) (*entities.Bill, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.BillUpdate) (*entities.Bill, error)
		DeleteOne(ctx context.Context, id uuid.UUID) error
		CheckExists(ctx context.Context, postId uuid.UUID) (bool, error)
		GetMonthlyRevenue(ctx context.Context, date time.Time) ([]string, []int64, error)
		GetRevenueForMonth(ctx context.Context, date time.Time) (int64, error)
		GetRevenueForDay(ctx context.Context, date time.Time) (int64, error)
	}
)
