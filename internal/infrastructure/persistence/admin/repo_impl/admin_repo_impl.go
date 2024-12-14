package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rAdmin struct {
	db *gorm.DB
}

func NewAdminRepositoryImplement(db *gorm.DB) *rAdmin {
	return &rAdmin{db: db}
}

func (r *rAdmin) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Admin, error) {
	return nil, nil
}

func (r *rAdmin) CreateOne(
	ctx context.Context,
	entity *entities.Admin,
) (*entities.Admin, error) {
	return nil, nil
}

func (r *rAdmin) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.AdminUpdate,
) (*entities.Admin, error) {
	return nil, nil
}

func (r *rAdmin) GetOne(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*entities.Admin, error) {
	return nil, nil
}

func (r *rAdmin) GetMany(
	ctx context.Context,
	query *query.GetManyAdminQuery,
) ([]*entities.Admin, *response.PagingResponse, error) {
	return nil, nil, nil
}

func (r *rAdmin) CheckAdminExistByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	return false, nil
}
