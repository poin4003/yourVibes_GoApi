package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rUserReport struct {
	db *gorm.DB
}

func NewUserReportRepositoryImplement(db *gorm.DB) *rUserReport {
	return &rUserReport{db: db}
}

func (r *rUserReport) GetByUserIdAndReportedUserId(
	ctx context.Context,
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) (*entities.UserReport, error) {
	return nil, nil
}

func (r *rUserReport) CreateOne(
	ctx context.Context,
	entity *entities.UserReport,
) (*entities.UserReport, error) {
	return nil, nil
}

func (r *rUserReport) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.UserReportUpdate,
) (*entities.UserReport, error) {
	return nil, nil
}

func (r *rUserReport) DeleteOne(
	ctx context.Context,
	id uuid.UUID,
) error {
	return nil
}

func (r *rUserReport) GetMany(
	ctx context.Context,
	query *query.GetManyUserReportQuery,
) ([]*entities.UserReport, *response.PagingResponse, error) {
	return nil, nil, nil
}
