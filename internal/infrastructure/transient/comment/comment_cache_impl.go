package comment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type tComment struct {
	client *redis.Client
}

func NewCommentCacheImplement(client *redis.Client) *tComment {
	return &tComment{client: client}
}

func (t *tComment) SetComment(
	ctx context.Context,
	comment *entities.Comment,
) {
	key := fmt.Sprintf("comment:%s", comment.ID.String())
	data, err := json.Marshal(comment)
	if err != nil {
		global.Logger.Warn("Failed to marshal comment", zap.String("comment_id", comment.ID.String()), zap.Error(err))
		return
	}

	if err = t.client.Set(ctx, key, string(data), consts.TTL_COMMON).Err(); err != nil {
		global.Logger.Warn("Failed to set comment to redis", zap.String("comment_id", comment.ID.String()), zap.Error(err))
	}

	userSetKey := fmt.Sprintf("comment_ids_by_user:%s", comment.UserId.String())
	if err = t.client.SAdd(ctx, userSetKey, comment.ID.String()).Err(); err != nil {
		global.Logger.Warn("Failed to add commentID to user set", zap.String("user_id", comment.UserId.String()), zap.Error(err))
	}
}

func (t *tComment) GetComment(
	ctx context.Context,
	commentID uuid.UUID,
) *entities.Comment {
	key := fmt.Sprintf("comment:%s", commentID.String())
	data, err := t.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		global.Logger.Warn("Failed to get comment", zap.String("comment_id", commentID.String()), zap.Error(err))
		return nil
	}
	comment := &entities.Comment{}
	if err = json.Unmarshal(data, comment); err != nil {
		global.Logger.Warn("Failed to unmarshal comment", zap.String("comment_id", commentID.String()), zap.Error(err))
		return nil
	}
	return comment
}

func (t *tComment) DeleteComment(
	ctx context.Context,
	commentID uuid.UUID,
) {
	key := fmt.Sprintf("comment:%s", commentID.String())
	if err := t.client.Del(ctx, key).Err(); err != nil {
		global.Logger.Warn("Failed to delete comment", zap.String("comment_id", commentID.String()), zap.Error(err))
	}
}

func (t *tComment) SetPostComment(
	ctx context.Context,
	postID uuid.UUID, parentID uuid.UUID, commentIds []uuid.UUID,
	paging *response.PagingResponse,
) {
	var key, totalKey string
	if parentID != uuid.Nil {
		key = fmt.Sprintf("post_comment:%s:%s", postID.String(), parentID.String())
		totalKey = fmt.Sprintf("post_comment:%s:%s:total", postID.String(), parentID.String())
	} else {
		key = fmt.Sprintf("post_comment:%s", postID.String())
		totalKey = fmt.Sprintf("post_comment:%s:total", postID.String())
	}

	zMembers := make([]redis.Z, len(commentIds))
	for i, id := range commentIds {
		zMembers[i] = redis.Z{
			Score:  float64(len(commentIds) - 1),
			Member: id.String(),
		}
	}

	pipe := t.client.Pipeline()
	if len(zMembers) > 0 {
		pipe.ZAdd(ctx, key, zMembers...)
		pipe.Expire(ctx, key, consts.TTL_COMMON)
		pipe.Set(ctx, totalKey, paging.Total, consts.TTL_COMMON)
	} else {
		pipe.Del(ctx, key, totalKey)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		global.Logger.Warn("Failed to set post comment", zap.String("post_id", postID.String()), zap.Error(err))
	}
}

func (t *tComment) GetPostComment(
	ctx context.Context,
	postID uuid.UUID, parentID uuid.UUID, limit, page int,
) ([]uuid.UUID, *response.PagingResponse) {
	var key, totalKey string
	if parentID != uuid.Nil {
		key = fmt.Sprintf("post_comment:%s:%s", postID.String(), parentID.String())
		totalKey = fmt.Sprintf("post_comment:%s:%s:total", postID.String(), parentID.String())
	} else {
		key = fmt.Sprintf("post_comment:%s", postID.String())
		totalKey = fmt.Sprintf("post_comment:%s:total", postID.String())
	}

	start := int64((page - 1) * limit)
	stop := start + int64(limit) - 1

	idStrings, err := t.client.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		global.Logger.Warn("Failed to get post comment", zap.String("post_id", postID.String()), zap.Error(err))
		return nil, nil
	}

	totalStr, err := t.client.Get(ctx, totalKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		global.Logger.Warn("Failed to get total post comment", zap.String("post_id", postID.String()), zap.Error(err))
		return nil, nil
	}

	total, err := strconv.ParseInt(totalStr, 10, 64)
	if err != nil {
		global.Logger.Warn("Failed to parse total post comment", zap.String("post_id", postID.String()), zap.Error(err))
		return nil, nil
	}

	var commentUUIDs []uuid.UUID
	for _, str := range idStrings {
		var id uuid.UUID
		if id, err = uuid.Parse(str); err == nil {
			commentUUIDs = append(commentUUIDs, id)
		} else {
			global.Logger.Warn("Failed to parse comment id", zap.String("comment_id", str), zap.Error(err))
		}
	}

	paging := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	return commentUUIDs, paging
}

func (t *tComment) DeletePostComment(ctx context.Context, postID uuid.UUID) {
	mainKey := fmt.Sprintf("post_comment:%s", postID.String())
	mainTotalKey := fmt.Sprintf("post_comment:%s:total", postID.String())
	childPattern := fmt.Sprintf("post_comment:%s:*", postID.String())

	var cursor uint64
	var keysToDelete []string
	for {
		keys, nextCursor, err := t.client.Scan(ctx, cursor, childPattern, 10).Result()
		if err != nil {
			global.Logger.Warn("Failed to get post comment", zap.String("post_id", postID.String()), zap.Error(err))
			return
		}
		keysToDelete = append(keysToDelete, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	keysToDelete = append(keysToDelete, mainKey, mainTotalKey)

	if len(keysToDelete) > 0 {
		if err := t.client.Del(ctx, keysToDelete...).Err(); err != nil {
			global.Logger.Warn("Failed to delete post comment keys", zap.String("post_id", postID.String()), zap.Error(err))
		}
	}
}

func (t *tComment) DeleteAllUserComments(ctx context.Context, userID uuid.UUID) {
	userSetKey := fmt.Sprintf("comment_ids_by_user:%s", userID.String())

	commentIDs, err := t.client.SMembers(ctx, userSetKey).Result()
	if err != nil {
		global.Logger.Error("Failed to get commentIDs of user", zap.String("user_id", userID.String()), zap.Error(err))
		return
	}

	if len(commentIDs) == 0 {
		return
	}

	var keys []string
	for _, id := range commentIDs {
		keys = append(keys, fmt.Sprintf("comment:%s", id))
	}
	keys = append(keys, userSetKey)

	if err = t.client.Del(ctx, keys...).Err(); err != nil {
		global.Logger.Error("Failed to delete user comments from redis", zap.String("user_id", userID.String()), zap.Error(err))
	}
}

func (t *tComment) DeleteAllCommentCache(ctx context.Context) error {
	patterns := []string{
		"comment:*",
		"post_comment:*",
		"personal_post:*",
		"comment_ids_by_user:*",
	}

	for _, pattern := range patterns {
		iter := t.client.Scan(ctx, 0, pattern, 0).Iterator()
		var keysToDelete []string

		for iter.Next(ctx) {
			keysToDelete = append(keysToDelete, iter.Val())
		}

		if err := iter.Err(); err != nil {
			return response.NewServerFailedError(err.Error())
		}

		if len(keysToDelete) > 0 {
			if err := t.client.Del(ctx, keysToDelete...).Err(); err != nil {
				return response.NewServerFailedError(err.Error())
			}
		}
	}
	return nil
}
