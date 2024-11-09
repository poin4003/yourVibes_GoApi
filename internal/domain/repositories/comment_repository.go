package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/comment/comment_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	ICommentRepository interface {
		CreateComment(ctx context.Context, comment *entities2.Comment) (*entities2.Comment, error)
		UpdateOneComment(ctx context.Context, commentId uuid.UUID, updateData map[string]interface{}) (*entities2.Comment, error)
		UpdateManyComment(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteOneComment(ctx context.Context, commentId uuid.UUID) (*entities2.Comment, error)
		DeleteManyComment(ctx context.Context, condition map[string]interface{}) error
		GetOneComment(ctx context.Context, query interface{}, args ...interface{}) (*entities2.Comment, error)
		GetManyComment(ctx context.Context, query *query.CommentQueryObject) ([]*entities2.Comment, *response.PagingResponse, error)
		GetMaxCommentRightByPostId(ctx context.Context, postId uuid.UUID) (int, error)
	}
	ILikeUserCommentRepository interface {
		CreateLikeUserComment(ctx context.Context, likeUserComment *entities2.LikeUserComment) error
		DeleteLikeUserComment(ctx context.Context, likeUserComment *entities2.LikeUserComment) error
		GetLikeUserComment(ctx context.Context, commentId uuid.UUID, query *query.CommentLikeQueryObject) ([]*entities2.User, *response.PagingResponse, error)
		CheckUserLikeComment(ctx context.Context, likeUserComment *entities2.LikeUserComment) (bool, error)
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
