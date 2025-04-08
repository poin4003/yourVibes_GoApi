package cache

import (
	"context"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type (
	IPostCache interface {
		SetPost(ctx context.Context, post *entities.Post)
		GetPost(ctx context.Context, postID uuid.UUID) *entities.Post
		DeletePost(ctx context.Context, postID uuid.UUID)
		SetFeeds(ctx context.Context, inputKey consts.RedisKey, userID uuid.UUID, postsIds []uuid.UUID, paging *response.PagingResponse)
		GetFeeds(ctx context.Context, inputKey consts.RedisKey, userID uuid.UUID, limit, page int) ([]uuid.UUID, *response.PagingResponse)
		DeleteFeeds(ctx context.Context, inputKey consts.RedisKey, userID uuid.UUID)
		DeleteFriendFeeds(ctx context.Context, inputKey consts.RedisKey, friendIDs []uuid.UUID)
	}
)
