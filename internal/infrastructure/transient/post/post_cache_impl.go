package post

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
	userID uuid.UUID, postIds []uuid.UUID, paging *response.PagingResponse,
) {
	key := fmt.Sprintf("%s:%s", inputKey, userID.String())
	pagingKey := fmt.Sprintf("%s:%s:paging", inputKey, userID.String())

	// Convert postIds to string
	postIdStrings := make([]interface{}, len(postIds))
	for i, id := range postIds {
		postIdStrings[i] = id.String()
	}
	// Save postIds into list
	pipe := t.client.Pipeline()
	pipe.RPush(ctx, key, postIdStrings...)
	pipe.Expire(ctx, key, consts.TTL_COMMON)

	// Save paging data into redis
	pagingData, err := json.Marshal(paging)
	if err != nil {
		global.Logger.Warn("Failed to marshal paging", zap.String("user_id", userID.String()), zap.Error(err))
		return
	}
	pipe.Set(ctx, pagingKey, string(pagingData), consts.TTL_COMMON)

	// Execute
	_, err = pipe.Exec(ctx)
	if err != nil {
		global.Logger.Warn("Failed to set personal posts to redis", zap.String("user_id", userID.String()), zap.Error(err))
	}
}

func (t *tPost) GetFeeds(
	ctx context.Context,
	inputKey consts.RedisKey,
	userID uuid.UUID, limit, page int,
) ([]uuid.UUID, *response.PagingResponse) {
	key := fmt.Sprintf("%s:%s", inputKey, userID.String())
	pagingKey := fmt.Sprintf("%s:%s:paging", inputKey, userID.String())

	pagingData, err := t.client.Get(ctx, pagingKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		global.Logger.Warn("Failed to get paging from redis", zap.String("user_id", userID.String()), zap.String("paging_key", pagingKey))
		return nil, nil
	}

	var paging response.PagingResponse
	if err = json.Unmarshal(pagingData, &paging); err != nil {
		global.Logger.Warn("Failed to unmarshal paging from redis", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, nil
	}

	offset := (page - 1) * limit
	postIds, err := t.client.LRange(ctx, key, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		global.Logger.Warn("Failed to get personal posts from redis", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, nil
	}

	var postUUIDs []uuid.UUID
	for _, postIdString := range postIds {
		var postID uuid.UUID
		if postID, err = uuid.Parse(postIdString); err != nil {
			global.Logger.Warn("Failed to parse post id", zap.String("post_id", postIdString))
		} else {
			postUUIDs = append(postUUIDs, postID)
		}
	}

	return postUUIDs, &paging
}

func (t *tPost) DeleteFeeds(ctx context.Context, inputKey consts.RedisKey, userID uuid.UUID) {
	key := fmt.Sprintf("%s:%s", inputKey, userID.String())
	pagingKey := fmt.Sprintf("%s:%s:paging", inputKey, userID.String())

	if err := t.client.Del(ctx, key, pagingKey).Err(); err != nil {
		global.Logger.Warn("Failed to delete feeds from redis", zap.String("user_id", userID.String()), zap.Error(err))
	}
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

			key := fmt.Sprintf("%s:%s", inputKey, userID.String())
			pagingKey := fmt.Sprintf("%s:%s:paging", inputKey, userID.String())
			if err := t.client.Del(ctx, key, pagingKey).Err(); err != nil {
				global.Logger.Warn("Failed to delete feeds from redis", zap.String("user_id", userID.String()), zap.Error(err))
			}
		}(id)
	}
	wg.Wait()
}
