package repositories

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/query"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
)

type (
	IReportRepository interface {
		GetUserReportById(ctx context.Context, reportId uuid.UUID) (*entities.UserReportEntity, error)
		GetPostReportById(ctx context.Context, reportId uuid.UUID) (*entities.PostReportEntity, error)
		GetCommentReportById(ctx context.Context, reportId uuid.UUID) (*entities.CommentReportEntity, error)
		CreatePostReport(ctx context.Context, entity *entities.PostReportEntity) error
		CreateUserReport(ctx context.Context, entity *entities.UserReportEntity) error
		CreateCommentReport(ctx context.Context, entity *entities.CommentReportEntity) error
		HandleUserReport(ctx context.Context, reportId, adminId uuid.UUID) error
		HandlePostReport(ctx context.Context, reportId, adminId uuid.UUID) error
		HandleCommentReport(ctx context.Context, reportId, adminId uuid.UUID) error
		ActivateUser(ctx context.Context, reportId uuid.UUID) error
		ActivatePost(ctx context.Context, reportId uuid.UUID) error
		ActivateComment(ctx context.Context, reportId uuid.UUID) error
		GetManyUserReport(ctx context.Context, query *query.GetManyReportQuery) ([]*entities.UserReportEntity, *response.PagingResponse, error)
		GetManyPostReport(ctx context.Context, query *query.GetManyReportQuery) ([]*entities.PostReportEntity, *response.PagingResponse, error)
		GetManyCommentReport(ctx context.Context, query *query.GetManyReportQuery) ([]*entities.CommentReportEntity, *response.PagingResponse, error)
		DeleteReportById(ctx context.Context, reportId uuid.UUID) error
	}
)

var (
	localReport IReportRepository
)

func Report() IReportRepository {
	if localReport == nil {
		panic("repository_implement localReport not found for interface IReport")
	}

	return localReport
}

func InitReportRepository(i IReportRepository) {
	localReport = i
}
