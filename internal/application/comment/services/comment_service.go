package services

import (
	"context"

	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
)

type (
	ICommentUser interface {
		CreateComment(ctx context.Context, command *command.CreateCommentCommand) (result *command.CreateCommentResult, err error)
		UpdateComment(ctx context.Context, command *command.UpdateCommentCommand) (result *command.UpdateCommentResult, err error)
		DeleteComment(ctx context.Context, command *command.DeleteCommentCommand) error
		GetManyComments(ctx context.Context, query *query.GetManyCommentQuery) (result *query.GetManyCommentsResult, err error)
	}
	ICommentLike interface {
		LikeComment(ctx context.Context, command *command.LikeCommentCommand) (result *command.LikeCommentResult, err error)
		GetUsersOnLikeComment(ctx context.Context, query *query.GetCommentLikeQuery) (result *query.GetCommentLikeResult, err error)
	}
	ICommentReport interface {
		CreateCommentReport(ctx context.Context, command *command.CreateReportCommentCommand) (result *command.CreateReportCommentCommandResult, err error)
		HandleCommentReport(ctx context.Context, command *command.HandleCommentReportCommand) error
		DeleteCommentReport(ctx context.Context, command *command.DeleteCommentReportCommand) error
		ActivateComment(ctx context.Context, command *command.ActivateCommentCommand) error
		GetDetailCommentReport(ctx context.Context, query *query.GetOneCommentReportQuery) (result *query.CommentReportQueryResult, err error)
		GetManyCommentReport(ctx context.Context, query *query.GetManyCommentReportQuery) (result *query.CommentReportQueryListResult, err error)
	}
)

var (
	localCommentUser   ICommentUser
	localCommentLike   ICommentLike
	localCommentReport ICommentReport
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

func CommentReport() ICommentReport {
	if localCommentReport == nil {
		panic("repository_implement localCommentReport not found for interface ICommentReport")
	}

	return localCommentReport
}

func InitCommentUser(i ICommentUser) {
	localCommentUser = i
}

func InitCommentLike(i ICommentLike) {
	localCommentLike = i
}

func InitCommentReport(i ICommentReport) {
	localCommentReport = i
}
