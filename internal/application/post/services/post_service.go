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
	}
	IPostLike interface {
		LikePost(ctx context.Context, command *command.LikePostCommand) (result *command.LikePostCommandResult, err error)
		GetUsersOnLikes(ctx context.Context, query *query.GetPostLikeQuery) (result *query.GetPostLikeQueryResult, err error)
	}
	//IPostShare interface {
	//	SharePost(ctx context.Context, postId uuid.UUID, userId uuid.UUID, shareInput *request.SharePostInput) (post *models.Post, resultCode int, httpStatusCode int, err error)
	//}
	//IPostNewFeed interface {
	//	DeleteNewFeed(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (resultCode int, httpStatusCode int, err error)
	//	GetNewFeeds(ctx context.Context, userId uuid.UUID, query *query.NewFeedQueryObject) (postDtos []*dto_response.PostDto, pagingResponse *response.PagingResponse, resultCode int, httpStatusCode int, err error)
	//}
)

var (
	localPostUser     IPostUser
	localLikeUserPost IPostLike
	//localPostShare    IPostShare
	//localPostNewFeed  IPostNewFeed
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

//func PostShare() IPostShare {
//	if localPostShare == nil {
//		panic("repository_implement localPostShare not found for interface IPostShare")
//	}
//
//	return localPostShare
//}
//
//func PostNewFeed() IPostNewFeed {
//	if localPostNewFeed == nil {
//		panic("repository_implement localPostNewFeed not found for interface IPostNewFeed")
//	}
//
//	return localPostNewFeed
//}

func InitPostUser(i IPostUser) {
	localPostUser = i
}

func InitLikeUserPost(i IPostLike) {
	localLikeUserPost = i
}

//func InitPostShare(i IPostShare) {
//	localPostShare = i
//}
//
//func InitPostNewFeed(i IPostNewFeed) {
//	localPostNewFeed = i
//}
