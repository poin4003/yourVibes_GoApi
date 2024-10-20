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
		CreateComment(ctx context.Context, commentModel *model.Comment) (comment *model.Comment, resultCode int, err error)
		UpdateComment(ctx context.Context, commentId uuid.UUID, updateData map[string]interface{}) (comment *model.Comment, resultCode int, err error)
		DeleteComment(ctx context.Context, commentId uuid.UUID) (resultCode int, err error)
		GetManyComments(ctx context.Context, query *query_object.CommentQueryObject) (comments []*model.Comment, resultCode int, pagingResponse *response.PagingResponse, err error)
	}
)

var (
	localCommentUser ICommentUser
)

func CommentUser() ICommentUser {
	if localCommentUser == nil {
		panic("repository_implement localCommentUser not found for interface ICommentUser")
	}

	return localCommentUser
}

func InitCommentUser(i ICommentUser) {
	localCommentUser = i
}
