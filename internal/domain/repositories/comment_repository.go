package repositories

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
)

type (
	ICommentRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Comment, error)
		CreateOne(ctx context.Context, entity *entities.Comment) (*entities.Comment, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.CommentUpdate) (*entities.Comment, error)
		UpdateMany(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteOne(ctx context.Context, id uuid.UUID) (*entities.Comment, error)
		DeleteMany(ctx context.Context, condition map[string]interface{}) (int64, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Comment, error)
		GetMany(ctx context.Context, query *query.GetManyCommentQuery) ([]*entities.Comment, *response.PagingResponse, error)
		DeleteCommentAndChildComment(ctx context.Context, commentId uuid.UUID) error
	}
	ILikeUserCommentRepository interface {
		CreateLikeUserComment(ctx context.Context, entity *entities.LikeUserComment) error
		DeleteLikeUserComment(ctx context.Context, entity *entities.LikeUserComment) error
		GetLikeUserComment(ctx context.Context, query *query.GetCommentLikeQuery) ([]*entities.User, *response.PagingResponse, error)
		CheckUserLikeComment(ctx context.Context, entity *entities.LikeUserComment) (bool, error)
		CheckUserLikeManyComment(ctx context.Context, query *query.CheckUserLikeManyCommentQuery) (map[uuid.UUID]bool, error)
	}
)
