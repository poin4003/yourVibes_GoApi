package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	ICommentRepository interface {
		CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
		UpdateOneComment(ctx context.Context, commentId uuid.UUID, updateData map[string]interface{}) (*models.Comment, error)
		UpdateManyComment(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteOneComment(ctx context.Context, commentId uuid.UUID) (*models.Comment, error)
		DeleteManyComment(ctx context.Context, condition map[string]interface{}) error
		GetOneComment(ctx context.Context, query interface{}, args ...interface{}) (*models.Comment, error)
		GetManyComment(ctx context.Context, query *query.CommentQueryObject) ([]*models.Comment, *response.PagingResponse, error)
		GetMaxCommentRightByPostId(ctx context.Context, postId uuid.UUID) (int, error)
	}
	ILikeUserCommentRepository interface {
		CreateLikeUserComment(ctx context.Context, likeUserComment *models.LikeUserComment) error
		DeleteLikeUserComment(ctx context.Context, likeUserComment *models.LikeUserComment) error
		GetLikeUserComment(ctx context.Context, commentId uuid.UUID, query *query.CommentLikeQueryObject) ([]*models.User, *response.PagingResponse, error)
		CheckUserLikeComment(ctx context.Context, likeUserComment *models.LikeUserComment) (bool, error)
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
