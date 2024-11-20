package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	ICommentRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Comment, error)
		CreateOne(ctx context.Context, entity *entities.Comment) (*entities.Comment, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.CommentUpdate) (*entities.Comment, error)
		UpdateMany(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteOne(ctx context.Context, id uuid.UUID) (*entities.Comment, error)
		DeleteMany(ctx context.Context, condition map[string]interface{}) error
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Comment, error)
		GetMany(ctx context.Context, query *query.GetManyCommentQuery) ([]*entities.Comment, *response.PagingResponse, error)
		GetMaxCommentRightByPostId(ctx context.Context, postId uuid.UUID) (int, error)
	}
	ILikeUserCommentRepository interface {
		CreateLikeUserComment(ctx context.Context, entity *entities.LikeUserComment) error
		DeleteLikeUserComment(ctx context.Context, entity *entities.LikeUserComment) error
		GetLikeUserComment(ctx context.Context, query *query.GetCommentLikeQuery) ([]*entities.User, *response.PagingResponse, error)
		CheckUserLikeComment(ctx context.Context, entity *entities.LikeUserComment) (bool, error)
	}
)

var (
	localComment         ICommentRepository
	localLikeUserComment ILikeUserCommentRepository
)

func Comment() ICommentRepository {
	if localComment == nil {
		panic("repository_implement localComment not found for interface IComment")
	}

	return localComment
}

func LikeUserComment() ILikeUserCommentRepository {
	if localLikeUserComment == nil {
		panic("repository_implement localLikeUserComment not found for interface ILikeUserComment")
	}

	return localLikeUserComment
}

func InitCommentRepository(i ICommentRepository) {
	localComment = i
}

func InitLikeUserCommentRepository(i ILikeUserCommentRepository) {
	localLikeUserComment = i
}
