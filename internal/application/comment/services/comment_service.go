package services

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
)

// jfl;kasjdfkl;asjdf;klsdjaf;lkajsdf;lksadj

type (
	ICommentUser interface {
		CreateComment(ctx context.Context, command *command.CreateCommentCommand) (result *command.CreateCommentResult, err error)
		UpdateComment(ctx context.Context, command *command.UpdateCommentCommand) (result *command.UpdateCommentResult, err error)
		DeleteComment(ctx context.Context, command *command.DeleteCommentCommand) (result *command.DeleteCommentResult, err error)
		GetManyComments(ctx context.Context, query *query.GetManyCommentQuery) (result *query.GetManyCommentsResult, err error)
	}
	ICommentLike interface {
		LikeComment(ctx context.Context, command *command.LikeCommentCommand) (result *command.LikeCommentResult, err error)
		GetUsersOnLikeComment(ctx context.Context, query *query.GetCommentLikeQuery) (result *query.GetCommentLikeResult, err error)
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
