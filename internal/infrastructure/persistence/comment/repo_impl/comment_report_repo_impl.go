package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rCommentReport struct {
	db *gorm.DB
}

func NewCommentReportRepositoryImplement(db *gorm.DB) *rCommentReport {
	return &rCommentReport{db: db}
}

func (r *rCommentReport) GetByUserIdAndReportedCommentId(
	ctx context.Context,
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) (*entities.CommentReport, error) {
	return nil, nil
}

func (r *rCommentReport) CreateOne(
	ctx context.Context,
	entity *entities.CommentReport,
) (*entities.CommentReport, error) {
	return nil, nil
}

func (r *rCommentReport) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.CommentReportUpdate,
) (*entities.CommentReport, error) {
	return nil, nil
}

func (r *rCommentReport) DeleteOne(
	ctx context.Context,
	id uuid.UUID,
) error {
	return nil
}

func (r *rCommentReport) GetMany(
	ctx context.Context,
	query *query.GetManyCommentReportQuery,
) ([]*entities.CommentReport, *response.PagingResponse, error) {
	return nil, nil, nil
}
