package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
)

type (
	IPostRepository interface {
		CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
		UpdatePost(ctx context.Context, postId uuid.UUID, updateData map[string]interface{}) (*model.Post, error)
		DeletePost(ctx context.Context, postId uuid.UUID) (*model.Post, error)
		GetPost(ctx context.Context, query interface{}, args ...interface{}) (*model.Post, error)
		GetManyPost(ctx context.Context, query *query_object.PostQueryObject) ([]*model.Post, error)
	}
	IMediaRepository interface {
		CreateMedia(ctx context.Context, media *model.Media) (*model.Media, error)
		UpdateMedia(ctx context.Context, mediaId uint, updateData map[string]interface{}) (*model.Media, error)
		DeleteMedia(ctx context.Context, mediaId uint) error
		GetMedia(ctx context.Context, query interface{}, args ...interface{}) (*model.Media, error)
		GetManyMedia(ctx context.Context, query interface{}, args ...interface{}) ([]*model.Media, error)
	}
	ILikeUserPostRepository interface {
		CreateLikeUserPost(ctx context.Context, likeUserPost *model.LikeUserPost) error
		DeleteLikeUserPost(ctx context.Context, likeUserPost *model.LikeUserPost) error
		GetLikeUserPost(ctx context.Context, postId uuid.UUID) ([]*model.User, error)
	}
)

var (
	localMedia        IMediaRepository
	localPost         IPostRepository
	localLikeUserPost ILikeUserPostRepository
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

func InitPostRepository(i IPostRepository) {
	localPost = i
}

func InitMediaRepository(i IMediaRepository) {
	localMedia = i
}

func InitLikeUserPostRepository(i ILikeUserPostRepository) {
	localLikeUserPost = i
}
