package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	ICommentUser interface {
		CreateComment(ctx context.Context, commentModel *models.Comment) (comment *models.Comment, resultCode int, httpStatusCode int, err error)
		UpdateComment(ctx context.Context, commentId uuid.UUID, updateData map[string]interface{}) (comment *models.Comment, resultCode int, httpStatusCode int, err error)
		DeleteComment(ctx context.Context, commentId uuid.UUID) (resultCode int, httpStatusCode int, err error)
		GetManyComments(ctx context.Context, query *query.CommentQueryObject, userId uuid.UUID) (commentDtos []*response.CommentDto, resultCode int, httpStatusCode int, pagingResponse *pkg_response.PagingResponse, err error)
	}
	ICommentLike interface {
		LikeComment(ctx context.Context, likeUserComment *models.LikeUserComment, userId uuid.UUID) (commentDto *response.CommentDto, resultCode int, httpStatusCode int, err error)
		GetUsersOnLikeComment(ctx context.Context, commentId uuid.UUID, query *query.CommentLikeQueryObject) (users []*models.User, resultCode int, httpStatusCode int, pagingResponse *pkg_response.PagingResponse, err error)
	}
)

var (
	localCommentUser ICommentUser
	localCommentLike ICommentLike
)

func CommentUser() ICommentUser {
	if localCommentUser == nil {
		panic("repository_implement localCommentUser not found for interface ICommentUser")
	}

	return localCommentUser
}

func CommentLike() ICommentLike {
	if localCommentLike == nil {
		panic("repository_implement localCommentLike not found for interface ICommentLike")
	}

	return localCommentLike
}

func InitCommentUser(i ICommentUser) {
	localCommentUser = i
}

func InitCommentLike(i ICommentLike) {
	localCommentLike = i
}
