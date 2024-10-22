package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	ICommentUser interface {
		CreateComment(ctx context.Context, commentModel *model.Comment) (comment *model.Comment, resultCode int, httpStatusCode int, err error)
		UpdateComment(ctx context.Context, commentId uuid.UUID, updateData map[string]interface{}) (comment *model.Comment, resultCode int, httpStatusCode int, err error)
		DeleteComment(ctx context.Context, commentId uuid.UUID) (resultCode int, httpStatusCode int, err error)
		GetManyComments(ctx context.Context, query *query_object.CommentQueryObject) (comments []*model.Comment, resultCode int, httpStatusCode int, pagingResponse *response.PagingResponse, err error)
	}
	ICommentLike interface {
		LikeComment(ctx context.Context, likeUserComment *model.LikeUserComment) (resultCode int, httpStatusCode int, err error)
		GetUsersOnLikeComment(ctx context.Context, commentId uuid.UUID, query *query_object.CommentLikeQueryObject) (users []*model.User, resultCode int, httpStatusCode int, pagingResponse *response.PagingResponse, err error)
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
