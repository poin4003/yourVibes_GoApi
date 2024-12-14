package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rPostReport struct {
	db *gorm.DB
}

func NewPostReportRepositoryImplement(db *gorm.DB) *rPostReport {
	return &rPostReport{db: db}
}

func (r *rPostReport) GetByUserIdAndReportedPostId(
	ctx context.Context,
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) (*entities.PostReport, error) {
	return nil, nil
}

func (r *rPostReport) CreateOne(
	ctx context.Context,
	entity *entities.PostReport,
) (*entities.PostReport, error) {
	return nil, nil
}

func (r *rPostReport) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.PostReportUpdate,
) (*entities.PostReport, error) {
	return nil, nil
}

func (r *rPostReport) DeleteOne(
	ctx context.Context,
	id uuid.UUID,
) error {
	return nil
}

func (r *rPostReport) GetMany(
	ctx context.Context,
	query *query.GetManyPostReportQuery,
) ([]*entities.PostReport, *response.PagingResponse, error) {
	return nil, nil, nil
}
