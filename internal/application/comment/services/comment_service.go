package services

import (
	"context"
	"github.com/google/uuid"
	entities2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/entities"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/comment/comment_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/comment/comment_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	ICommentUser interface {
		CreateComment(ctx context.Context, commentModel *entities2.Comment) (comment *entities2.Comment, resultCode int, httpStatusCode int, err error)
		UpdateComment(ctx context.Context, commentId uuid.UUID, updateData map[string]interface{}) (comment *entities2.Comment, resultCode int, httpStatusCode int, err error)
		DeleteComment(ctx context.Context, commentId uuid.UUID) (resultCode int, httpStatusCode int, err error)
		GetManyComments(ctx context.Context, query *query.CommentQueryObject, userId uuid.UUID) (commentDtos []*response2.CommentDto, resultCode int, httpStatusCode int, pagingResponse *response.PagingResponse, err error)
	}
	ICommentLike interface {
		LikeComment(ctx context.Context, likeUserComment *entities2.LikeUserComment, userId uuid.UUID) (commentDto *response2.CommentDto, resultCode int, httpStatusCode int, err error)
		GetUsersOnLikeComment(ctx context.Context, commentId uuid.UUID, query *query.CommentLikeQueryObject) (users []*entities2.User, resultCode int, httpStatusCode int, pagingResponse *response.PagingResponse, err error)
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
