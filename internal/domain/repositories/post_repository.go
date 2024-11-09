package repositories

import (
	"context"
	"github.com/google/uuid"
	entities2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/post/post_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	IPostRepository interface {
		CreatePost(ctx context.Context, post *entities2.Post) (*entities2.Post, error)
		UpdatePost(ctx context.Context, postId uuid.UUID, updateData map[string]interface{}) (*entities2.Post, error)
		DeletePost(ctx context.Context, postId uuid.UUID) (*entities2.Post, error)
		GetPost(ctx context.Context, query interface{}, args ...interface{}) (*entities2.Post, error)
		GetManyPost(ctx context.Context, query *query.PostQueryObject) ([]*entities2.Post, *response.PagingResponse, error)
	}
	IMediaRepository interface {
		CreateMedia(ctx context.Context, media *entities2.Media) (*entities2.Media, error)
		UpdateMedia(ctx context.Context, mediaId uint, updateData map[string]interface{}) (*entities2.Media, error)
		DeleteMedia(ctx context.Context, mediaId uint) error
		GetMedia(ctx context.Context, query interface{}, args ...interface{}) (*entities2.Media, error)
		GetManyMedia(ctx context.Context, query interface{}, args ...interface{}) ([]*entities2.Media, error)
	}
	ILikeUserPostRepository interface {
		CreateLikeUserPost(ctx context.Context, likeUserPost *entities2.LikeUserPost) error
		DeleteLikeUserPost(ctx context.Context, likeUserPost *entities2.LikeUserPost) error
		GetLikeUserPost(ctx context.Context, postId uuid.UUID, query *query.PostLikeQueryObject) ([]*entities2.User, *response.PagingResponse, error)
		CheckUserLikePost(ctx context.Context, likeUserPost *entities2.LikeUserPost) (bool, error)
	}
	INewFeedRepository interface {
		CreateManyNewFeed(ctx context.Context, postId uuid.UUID, friendIds []uuid.UUID) error
		DeleteNewFeed(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
		GetManyNewFeed(ctx context.Context, userId uuid.UUID, query *query.NewFeedQueryObject) ([]*entities2.Post, *response.PagingResponse, error)
	}
)

var (
	localMedia        IMediaRepository
	localPost         IPostRepository
	localLikeUserPost ILikeUserPostRepository
	localNewFeed      INewFeedRepository
)

func Post() IPostRepository {
	if localPost == nil {
		panic("repository_implement localPost not found for interface IPost")
	}

	return localPost
}

func Media() IMediaRepository {
	if localMedia == nil {
		panic("repository_implement localMedia not found for interface IMedia")
	}

	return localMedia
}

func LikeUserPost() ILikeUserPostRepository {
	if localLikeUserPost == nil {
		panic("repository_implement localLikeUserPost not found for interface ILikeUserPost")
	}

	return localLikeUserPost
}

func NewFeed() INewFeedRepository {
	if localNewFeed == nil {
		panic("repository_implement localNewFeed not found for interface INewFeed")
	}

	return localNewFeed
}

func InitPostRepository(i IPostRepository) {
	localPost = i
}

func InitMediaRepository(i IMediaRepository) {
	localMedia = i
}

func InitLikeUserPostRepository(i ILikeUserPostRepository) {
	localLikeUserPost = i
}

func InitNewFeedRepository(i INewFeedRepository) {
	localNewFeed = i
}
