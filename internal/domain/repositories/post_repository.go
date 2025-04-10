package repositories

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
)

type (
	IPostRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Post, error)
		CreateOne(ctx context.Context, entity *entities.Post) (*entities.Post, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.PostUpdate) (*entities.Post, error)
		UpdateMany(ctx context.Context, condition map[string]interface{}, updateData *entities.PostUpdate) error
		DeleteOne(ctx context.Context, id uuid.UUID) (*entities.Post, error)
		GetMany(ctx context.Context, query *query.GetManyPostQuery) ([]*entities.Post, *response.PagingResponse, error)
		GetTrendingPost(ctx context.Context, query *query.GetTrendingPostQuery) ([]*entities.Post, *response.PagingResponse, error)
		UpdateExpiredAdvertisements(ctx context.Context) error
		CheckPostOwner(ctx context.Context, postId uuid.UUID, userId uuid.UUID) (bool, error)
		GetTotalPostCount(ctx context.Context) (int, error)
		GetTotalPostCountByUserId(ctx context.Context, userId uuid.UUID) (int64, error)
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
		CheckUserLikeManyPost(ctx context.Context, query *query.CheckUserLikeManyPostQuery) (map[uuid.UUID]bool, error)
	}
	INewFeedRepository interface {
		CreateMany(ctx context.Context, postId uuid.UUID, userId uuid.UUID) error
		DeleteOne(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
		DeleteMany(ctx context.Context, condition map[string]interface{}) error
		GetMany(ctx context.Context, query *query.GetNewFeedQuery) ([]*entities.Post, *response.PagingResponse, error)
		CreateManyWithRandomUser(ctx context.Context, numUsers int) error
		DeleteExpiredAdvertiseFromNewFeeds(ctx context.Context) error
		CreateManyFeaturedPosts(ctx context.Context, numUsers int) error
	}
)
