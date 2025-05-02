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
		ClearAllCommentCaches(ctx context.Context) error
	}
	ICommentLike interface {
		LikeComment(ctx context.Context, command *command.LikeCommentCommand) (result *command.LikeCommentResult, err error)
		GetUsersOnLikeComment(ctx context.Context, query *query.GetCommentLikeQuery) (result *query.GetCommentLikeResult, err error)
	}
)
