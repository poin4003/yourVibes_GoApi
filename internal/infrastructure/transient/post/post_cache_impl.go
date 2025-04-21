package post

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

type tPost struct {
	client *redis.Client
}

func NewPostCacheImplement(client *redis.Client) *tPost {
	return &tPost{client: client}
}

func (t *tPost) SetPost(ctx context.Context, post *entities.Post) {
	key := fmt.Sprintf("post:%s", post.ID.String())
	data, err := json.Marshal(post)
	if err != nil {
		global.Logger.Warn("Failed to marshal post", zap.String("post_id", post.ID.String()), zap.Error(err))
		return
	}
	if err = t.client.Set(ctx, key, string(data), consts.TTL_COMMON).Err(); err != nil {
		global.Logger.Warn("Failed to set post to redis", zap.String("post_id", post.ID.String()), zap.Error(err))
	}
}

func (t *tPost) GetPost(ctx context.Context, postID uuid.UUID) *entities.Post {
	key := fmt.Sprintf("post:%s", postID.String())
	data, err := t.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		global.Logger.Warn("Failed to get post to redis", zap.Error(err))
		return nil
	}
	post := &entities.Post{}
	if err = json.Unmarshal(data, post); err != nil {
		global.Logger.Warn("Failed to unmarshal post from redis", zap.String("post_id", postID.String()), zap.Error(err))
		return nil
	}
	return post
}

func (t *tPost) DeletePost(ctx context.Context, postID uuid.UUID) {
	key := fmt.Sprintf("post:%s", postID.String())
	if err := t.client.Del(ctx, key).Err(); err != nil {
		global.Logger.Warn("Failed to delete post from redis", zap.Error(err))
	}
}

func (t *tPost) SetFeeds(
	ctx context.Context,
	inputKey consts.RedisKey,
	userID uuid.UUID,
	postIDs []uuid.UUID,
	paging *response.PagingResponse,
) {
	key := fmt.Sprintf("%s:%s", inputKey, userID.String())

	zMembers := make([]redis.Z, len(postIDs))
	for i, id := range postIDs {
		zMembers[i] = redis.Z{
			Score:  float64(len(postIDs) - i),
			Member: id.String(),
		}
	}

	if err := t.client.ZAdd(ctx, key, zMembers...).Err(); err != nil {
		global.Logger.Warn("Failed to set feeds", zap.String("key", key), zap.Error(err))
		return
	}

	totalKey := fmt.Sprintf("%s:total:%s", inputKey, userID.String())
	if err := t.client.Set(ctx, totalKey, paging.Total, 0).Err(); err != nil {
		global.Logger.Warn("Failed to set total paging", zap.String("key", totalKey), zap.Error(err))
	}
}

func (t *tPost) GetFeeds(
	ctx context.Context,
	inputKey consts.RedisKey,
	userID uuid.UUID, limit, page int,
) ([]uuid.UUID, *response.PagingResponse) {
	key := fmt.Sprintf("%s:%s", inputKey, userID.String())
	totalKey := fmt.Sprintf("%s:total:%s", inputKey, userID.String())

	start := int64((page - 1) * limit)
	stop := start + int64(limit) - 1

	idStrings, err := t.client.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		global.Logger.Warn("Failed to get post feeds", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, nil
	}

	totalStr, err := t.client.Get(ctx, totalKey).Result()
	if err != nil {
		global.Logger.Warn("Failed to get total from cache", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, nil
	}

	total, err := strconv.ParseInt(totalStr, 10, 64)
	if err != nil {
		global.Logger.Warn("Failed to parse total", zap.String("totalStr", totalStr), zap.Error(err))
		return nil, nil
	}

	var ids []uuid.UUID
	for _, str := range idStrings {
		var id uuid.UUID
		id, err = uuid.Parse(str)
		if err == nil {
			ids = append(ids, id)
		} else {
			global.Logger.Warn("Failed to parse id", zap.String("id", str), zap.Error(err))
		}
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	return ids, pagingResponse
}

func (t *tPost) DeleteFeeds(ctx context.Context, inputKey consts.RedisKey, userID uuid.UUID) {
	key := fmt.Sprintf("%s:%s", inputKey, userID.String())
	totalKey := fmt.Sprintf("%s:total:%s", inputKey, userID.String())

	t.client.Del(ctx, key)
	t.client.Del(ctx, totalKey)
}

func (t *tPost) DeleteFriendFeeds(ctx context.Context, inputKey consts.RedisKey, friendIDs []uuid.UUID) {
	var wg sync.WaitGroup
	maxWorkers := 10
	sem := make(chan struct{}, maxWorkers)

	for _, id := range friendIDs {
		wg.Add(1)
		sem <- struct{}{}

		go func(userID uuid.UUID) {
			defer wg.Done()
			defer func() { <-sem }()
			t.DeleteFeeds(ctx, inputKey, userID)
		}(id)
	}
	wg.Wait()
}

func (t *tPost) DeleteRelatedPost(
	ctx context.Context,
	inputKey consts.RedisKey,
	userID uuid.UUID,
) {
	listKey := fmt.Sprintf("%s:%s", inputKey, userID.String())
	totalKey := fmt.Sprintf("%s:total:%s", inputKey, userID.String())

	postIDs, err := t.client.ZRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		global.Logger.Warn("Failed to get related post from redis", zap.String("user_id", userID.String()), zap.Error(err))
		return
	}

	if len(postIDs) == 0 {
		return
	}

	var keysToDelete []string
	for _, id := range postIDs {
		keysToDelete = append(keysToDelete, fmt.Sprintf("post:%s", id))
	}

	keysToDelete = append(keysToDelete, listKey, totalKey)

	if err = t.client.Del(ctx, keysToDelete...).Err(); err != nil {
		global.Logger.Warn("Failed to delete post from redis", zap.String("user_id", userID.String()), zap.Error(err))
	}
}

func (t *tPost) SetPostForCreate(
	ctx context.Context,
	postID uuid.UUID,
	post *entities.PostForCreate,
) error {
	key := fmt.Sprintf("create_post:%s", postID.String())
	data, err := json.Marshal(&post)
	if err != nil {
		global.Logger.Warn("Failed to marshal create post", zap.String("post_id", postID.String()), zap.Error(err))
		return response.NewServerFailedError(err.Error())
	}
	if err = t.client.Set(ctx, key, string(data), consts.TTL_COMMON).Err(); err != nil {
		global.Logger.Warn("Failed to set create post to redis", zap.String("post_id", postID.String()), zap.Error(err))
		return response.NewServerFailedError(err.Error())
	}
	return nil
}

func (t *tPost) GetPostForCreate(
	ctx context.Context,
	postID uuid.UUID,
) (*entities.PostForCreate, error) {
	key := fmt.Sprintf("create_post:%s", postID.String())
	data, err := t.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, err
		}
		global.Logger.Warn("Failed to get post for create from redis", zap.Error(err))
		return nil, err
	}
	post := &entities.PostForCreate{}
	if err = json.Unmarshal(data, post); err != nil {
		global.Logger.Warn("Failed to unmarshal post for create", zap.String("post_id", postID.String()), zap.Error(err))
		return nil, err
	}
	return post, nil
}

func (t *tPost) DeletePostForCreate(
	ctx context.Context,
	postID uuid.UUID,
) error {
	key := fmt.Sprintf("create_post:%s", postID.String())
	if err := t.client.Del(ctx, key).Err(); err != nil {
		global.Logger.Warn("Failed to delete post for create from redis", zap.Error(err))
		return response.NewServerFailedError(err.Error())
	}
	return nil
}
