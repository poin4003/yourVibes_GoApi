package services

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
)

type (
	IPostUser interface {
		CreatePost(ctx context.Context, command *command.CreatePostCommand) (result *command.CreatePostCommandResult, err error)
		UpdatePost(ctx context.Context, command *command.UpdatePostCommand) (result *command.UpdatePostCommandResult, err error)
		DeletePost(ctx context.Context, command *command.DeletePostCommand) (result *command.DeletePostCommandResult, err error)
		GetPost(ctx context.Context, query *query.GetOnePostQuery) (result *query.GetOnePostQueryResult, err error)
		GetManyPosts(ctx context.Context, query *query.GetManyPostQuery) (result *query.GetManyPostQueryResult, err error)
		CheckPostOwner(ctx context.Context, query *query.CheckPostOwnerQuery) (bool, error)
	}
	IPostLike interface {
		LikePost(ctx context.Context, command *command.LikePostCommand) (result *command.LikePostCommandResult, err error)
		GetUsersOnLikes(ctx context.Context, query *query.GetPostLikeQuery) (result *query.GetPostLikeQueryResult, err error)
	}
	IPostShare interface {
		SharePost(ctx context.Context, command *command.SharePostCommand) (result *command.SharePostCommandResult, err error)
	}
	IPostNewFeed interface {
		DeleteNewFeed(ctx context.Context, command *command.DeleteNewFeedCommand) (result *command.DeleteNewFeedCommandResult, err error)
		GetNewFeeds(ctx context.Context, query *query.GetNewFeedQuery) (result *query.GetNewFeedResult, err error)
	}
	IPostReport interface {
		CreatePostReport(ctx context.Context, command *command.CreateReportPostCommand) (result *command.CreateReportPostCommandResult, err error)
		HandlePostReport(ctx context.Context, command *command.HandlePostReportCommand) (result *command.HandlePostReportCommandResult, err error)
		DeletePostReport(ctx context.Context, command *command.DeletePostReportCommand) (result *command.DeletePostReportCommandResult, err error)
		ActivatePost(ctx context.Context, command *command.ActivatePostCommand) (result *command.ActivatePostCommandResult, err error)
		GetDetailPostReport(ctx context.Context, query *query.GetOnePostReportQuery) (result *query.PostReportQueryResult, err error)
		GetManyPostReport(ctx context.Context, query *query.GetManyPostReportQuery) (result *query.PostReportQueryListResult, err error)
	}
)

var (
	localPostUser     IPostUser
	localLikeUserPost IPostLike
	localPostShare    IPostShare
	localPostNewFeed  IPostNewFeed
	localPostReport   IPostReport
)

func PostUser() IPostUser {
	if localPostUser == nil {
		panic("service_implement localPostUser not found for interface IPostUser")
	}

	return localPostUser
}

func LikeUserPost() IPostLike {
	if localLikeUserPost == nil {
		panic("repository_implement localLikeUserPost not found for interface ILikeUserPost")
	}

	return localLikeUserPost
}

func PostShare() IPostShare {
	if localPostShare == nil {
		panic("repository_implement localPostShare not found for interface IPostShare")
	}

	return localPostShare
}

func PostNewFeed() IPostNewFeed {
	if localPostNewFeed == nil {
		panic("repository_implement localPostNewFeed not found for interface IPostNewFeed")
	}

	return localPostNewFeed
}

func PostReport() IPostReport {
	if localPostReport == nil {
		panic("repository_implement localPostReport not found for interface IPostReport")
	}

	return localPostReport
}

func InitPostUser(i IPostUser) {
	localPostUser = i
}

func InitLikeUserPost(i IPostLike) {
	localLikeUserPost = i
}

func InitPostShare(i IPostShare) {
	localPostShare = i
}

func InitPostNewFeed(i IPostNewFeed) {
	localPostNewFeed = i
}

func InitPostReport(i IPostReport) {
	localPostReport = i
}
