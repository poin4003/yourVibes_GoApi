package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	dto_response "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"mime/multipart"
)

type (
	IPostUser interface {
		CreatePost(ctx context.Context, postModel *models.Post, inMedia []multipart.File) (post *models.Post, resultCode int, httpStatusCode int, err error)
		UpdatePost(ctx context.Context, postId uuid.UUID, updateData map[string]interface{}, deleteMediaIds []uint, inMedia []multipart.File) (post *models.Post, resultCode int, httpStatusCode int, err error)
		DeletePost(ctx context.Context, postId uuid.UUID) (resultCode int, httpStatusCode int, err error)
		GetPost(ctx context.Context, postId uuid.UUID, userId uuid.UUID) (postDto *dto_response.PostDto, resultCode int, httpStatusCode int, err error)
		GetManyPosts(ctx context.Context, query *query.PostQueryObject, userId uuid.UUID) (postDtos []*dto_response.PostDto, resultCode int, httpStatusCode int, response *response.PagingResponse, err error)
	}
	IPostLike interface {
		LikePost(ctx context.Context, likeUserPost *models.LikeUserPost, userId uuid.UUID) (postDto *dto_response.PostDto, resultCode int, httpStatusCode int, err error)
		GetUsersOnLikes(ctx context.Context, postId uuid.UUID, query *query.PostLikeQueryObject) (users []*models.User, resultCode int, httpStatusCode int, pagingResponse *response.PagingResponse, err error)
	}
	IPostShare interface {
		SharePost(ctx context.Context, postId uuid.UUID, userId uuid.UUID, shareInput *request.SharePostInput) (post *models.Post, resultCode int, httpStatusCode int, err error)
	}
	IPostNewFeed interface {
		DeleteNewFeed(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (resultCode int, httpStatusCode int, err error)
		GetNewFeeds(ctx context.Context, userId uuid.UUID, query *query.NewFeedQueryObject) (postDtos []*dto_response.PostDto, pagingResponse *response.PagingResponse, resultCode int, httpStatusCode int, err error)
	}
)

var (
	localPostUser     IPostUser
	localLikeUserPost IPostLike
	localPostShare    IPostShare
	localPostNewFeed  IPostNewFeed
)

func PostUser() IPostUser {
	if localPostUser == nil {
		panic("repository_implement localPostUser not found for interface IPostUser")
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
