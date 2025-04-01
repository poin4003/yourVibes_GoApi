package comment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	var comment *entities.Comment
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
	postID uuid.UUID, commentIds []uuid.UUID, paging *response.PagingResponse,
) {
	key := fmt.Sprintf("post_comment:%s", postID.String())
	pagingKey := fmt.Sprintf("post_comment:%s:paging", postID.String())

	// Convert commentIds to string
	commentIdString := make([]interface{}, len(commentIds))
	for i, id := range commentIds {
		commentIdString[i] = id.String()
	}
	// Save commentIds into list
	pipe := t.client.Pipeline()
	pipe.RPush(ctx, key, commentIdString...)
	pipe.Expire(ctx, pagingKey, consts.TTL_COMMON)

	// Save paging data into redis
	pagingData, err := json.Marshal(paging)
	if err != nil {
		global.Logger.Warn("Failed to marshal post comment", zap.String("post_id", postID.String()), zap.Error(err))
		return
	}
	pipe.Set(ctx, pagingKey, string(pagingData), consts.TTL_COMMON)

	// Execute
	_, err = pipe.Exec(ctx)
	if err != nil {
		global.Logger.Warn("Failed to set post comment", zap.String("post_id", postID.String()), zap.Error(err))
	}
}

func (t *tComment) GetPostComment(
	ctx context.Context,
	postID uuid.UUID, limit, page int,
) ([]uuid.UUID, *response.PagingResponse) {
	key := fmt.Sprintf("post_comment:%s", postID.String())
	pagingKey := fmt.Sprintf("post_comment:%s:paging", postID.String())

	pagingData, err := t.client.Get(ctx, pagingKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		global.Logger.Warn("Failed to get post comment", zap.String("post_id", postID.String()), zap.Error(err))
		return nil, nil
	}

	var paging response.PagingResponse
	if err = json.Unmarshal(pagingData, &paging); err != nil {
		global.Logger.Warn("Failed to unmarshal post comment", zap.String("post_id", postID.String()), zap.Error(err))
		return nil, nil
	}

	offset := (page - 1) * limit
	commentIds, err := t.client.LRange(ctx, key, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		global.Logger.Warn("Failed to get post comment", zap.String("post_id", postID.String()), zap.Error(err))
		return nil, nil
	}

	var commentUUIDs []uuid.UUID
	for _, commentString := range commentIds {
		var commentID uuid.UUID
		if commentID, err = uuid.Parse(commentString); err != nil {
			global.Logger.Warn("Failed to parse post comment", zap.String("post_id", postID.String()), zap.Error(err))
		} else {
			commentUUIDs = append(commentUUIDs, commentID)
		}
	}

	return commentUUIDs, &paging
}

func (t *tComment) DeletePostComment(ctx context.Context, postID uuid.UUID) {
	key := fmt.Sprintf("post_comment:%s", postID.String())
	pagingKey := fmt.Sprintf("post_comment:%s:paging", postID.String())

	if err := t.client.Del(ctx, key, pagingKey).Err(); err != nil {
		global.Logger.Warn("Failed to delete post comment", zap.String("post_id", postID.String()), zap.Error(err))
	}
}
