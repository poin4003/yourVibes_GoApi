package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type (
	IAdminRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Admin, error)
		GetStatusById(ctx context.Context, id uuid.UUID) (*bool, error)
		CreateOne(ctx context.Context, entity *entities.Admin) (*entities.Admin, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.AdminUpdate) (*entities.Admin, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Admin, error)
		GetMany(ctx context.Context, query *query.GetManyAdminQuery) ([]*entities.Admin, *response.PagingResponse, error)
		CheckAdminExistByEmail(ctx context.Context, email string) (bool, error)
	}
)

var (
	localAdmin IAdminRepository
)

func Admin() IAdminRepository {
	if localAdmin == nil {
		panic("repository_implement localAdmin not found for interface IAdminRepository")
	}

	return localAdmin
}

func InitAdminRepository(i IAdminRepository) {
	localAdmin = i
}
