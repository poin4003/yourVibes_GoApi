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
		HandleUserReport(ctx context.Context, reportId, adminId uuid.UUID) (*entities.UserForReport, error)
		HandlePostReport(ctx context.Context, reportId, adminId uuid.UUID) (*entities.PostForReport, error)
		HandleCommentReport(ctx context.Context, reportId, adminId uuid.UUID) (*entities.CommentForReport, error)
		ActivateUser(ctx context.Context, reportId uuid.UUID) (*entities.UserForReport, error)
		ActivatePost(ctx context.Context, reportId uuid.UUID) (*entities.PostForReport, error)
		ActivateComment(ctx context.Context, reportId uuid.UUID) (*entities.CommentForReport, error)
		GetManyUserReport(ctx context.Context, query *query.GetManyReportQuery) ([]*entities.UserReportEntity, *response.PagingResponse, error)
		GetManyPostReport(ctx context.Context, query *query.GetManyReportQuery) ([]*entities.PostReportEntity, *response.PagingResponse, error)
		GetManyCommentReport(ctx context.Context, query *query.GetManyReportQuery) ([]*entities.CommentReportEntity, *response.PagingResponse, error)
		DeleteReportById(ctx context.Context, reportId uuid.UUID) error
	}
)
