package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	IPostRepository interface {
		CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
		UpdatePost(ctx context.Context, postId uuid.UUID, updateData map[string]interface{}) (*models.Post, error)
		DeletePost(ctx context.Context, postId uuid.UUID) (*models.Post, error)
		GetPost(ctx context.Context, query interface{}, args ...interface{}) (*models.Post, error)
		GetManyPost(ctx context.Context, query *query.PostQueryObject) ([]*models.Post, *response.PagingResponse, error)
	}
	IMediaRepository interface {
		CreateMedia(ctx context.Context, media *models.Media) (*models.Media, error)
		UpdateMedia(ctx context.Context, mediaId uint, updateData map[string]interface{}) (*models.Media, error)
		DeleteMedia(ctx context.Context, mediaId uint) error
		GetMedia(ctx context.Context, query interface{}, args ...interface{}) (*models.Media, error)
		GetManyMedia(ctx context.Context, query interface{}, args ...interface{}) ([]*models.Media, error)
	}
	ILikeUserPostRepository interface {
		CreateLikeUserPost(ctx context.Context, likeUserPost *models.LikeUserPost) error
		DeleteLikeUserPost(ctx context.Context, likeUserPost *models.LikeUserPost) error
		GetLikeUserPost(ctx context.Context, postId uuid.UUID, query *query.PostLikeQueryObject) ([]*models.User, *response.PagingResponse, error)
		CheckUserLikePost(ctx context.Context, likeUserPost *models.LikeUserPost) (bool, error)
	}
	INewFeedRepository interface {
		CreateManyNewFeed(ctx context.Context, postId uuid.UUID, friendIds []uuid.UUID) error
		DeleteNewFeed(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
		GetManyNewFeed(ctx context.Context, userId uuid.UUID, query *query.NewFeedQueryObject) ([]*models.Post, *response.PagingResponse, error)
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
