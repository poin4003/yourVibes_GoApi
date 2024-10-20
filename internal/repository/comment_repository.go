package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	ICommentRepository interface {
		CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error)
		UpdateOneComment(ctx context.Context, commentId uuid.UUID, updateData map[string]interface{}) (*model.Comment, error)
		UpdateManyComment(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteComment(ctx context.Context, commentId uuid.UUID) (*model.Comment, error)
		GetComment(ctx context.Context, query interface{}, args ...interface{}) (*model.Comment, error)
		GetManyComment(ctx context.Context, query *query_object.CommentQueryObject) ([]*model.Comment, *response.PagingResponse, error)
		GetMaxCommentRightByPostId(ctx context.Context, postId uuid.UUID) (int, error)
	}
)

var (
	localComment ICommentRepository
)

func Comment() ICommentRepository {
	if localComment == nil {
		panic("repository_implement localComment not found for interface IComment")
	}

	return localComment
}

func InitCommentRepository(i ICommentRepository) {
	localComment = i
}
