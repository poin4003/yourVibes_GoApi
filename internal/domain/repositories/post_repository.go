package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	IPostRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Post, error)
		CreateOne(ctx context.Context, entity *entities.Post) (*entities.Post, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.PostUpdate) (*entities.Post, error)
		DeleteOne(ctx context.Context, id uuid.UUID) (*entities.Post, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Post, error)
		GetMany(ctx context.Context, query *query.GetManyPostQuery) ([]*entities.Post, *response.PagingResponse, error)
	}
	IMediaRepository interface {
		GetById(ctx context.Context, id uint) (*entities.Media, error)
		CreateOne(ctx context.Context, entity *entities.Media) (*entities.Media, error)
		UpdateOne(ctx context.Context, id uint, updateData *entities.MediaUpdate) (*entities.Media, error)
		DeleteOne(ctx context.Context, id uint) error
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Media, error)
		GetMany(ctx context.Context, query interface{}, args ...interface{}) ([]*entities.Media, error)
	}
	ILikeUserPostRepository interface {
		CreateLikeUserPost(ctx context.Context, entity *entities.LikeUserPost) error
		DeleteLikeUserPost(ctx context.Context, entity *entities.LikeUserPost) error
		GetLikeUserPost(ctx context.Context, query *query.GetPostLikeQuery) ([]*entities.User, *response.PagingResponse, error)
		CheckUserLikePost(ctx context.Context, entity *entities.LikeUserPost) (bool, error)
	}
	INewFeedRepository interface {
		CreateMany(ctx context.Context, postId uuid.UUID, friendIds []uuid.UUID) error
		DeleteOne(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
		GetMany(ctx context.Context, query *query.GetNewFeedQuery) ([]*entities.Post, *response.PagingResponse, error)
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
