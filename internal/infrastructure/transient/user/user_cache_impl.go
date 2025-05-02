package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type tUser struct {
	client *redis.Client
}

func NewUserCache(client *redis.Client) *tUser {
	return &tUser{client: client}
}

func (t *tUser) SetUserStatus(
	ctx context.Context, userId uuid.UUID, status bool,
) {
	key := fmt.Sprintf("user:%s:status", userId.String())
	val := "0"
	if status {
		val = "1"
	}
	err := t.client.Set(ctx, key, val, consts.TTL_COMMON).Err()
	if err != nil {
		global.Logger.Warn("Failed to set user status in redis", zap.String("user_id", userId.String()), zap.Error(err))
	}
}

func (t *tUser) GetUserStatus(
	ctx context.Context, userId uuid.UUID,
) *bool {
	key := fmt.Sprintf("user:%s:status", userId.String())
	val, err := t.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		global.Logger.Warn("Failed to get user status from redis", zap.Error(err))
		return nil
	}
	status := val == "1"
	return &status
}

func (t *tUser) DeleteUserStatus(
	ctx context.Context, userId uuid.UUID,
) {
	key := fmt.Sprintf("user:%s:status", userId.String())
	if err := t.client.Del(ctx, key).Err(); err != nil {
		global.Logger.Warn("Failed to delete user status from redis", zap.Error(err))
	}
}

func (t *tUser) SetOnline(
	ctx context.Context,
	userId uuid.UUID,
) {
	key := fmt.Sprintf("user:online:%s", userId.String())

	err := t.client.Set(ctx, key, "1", 5*time.Minute).Err()
	if err != nil {
		global.Logger.Warn("Failed to set online in redis", zap.Error(err))
	}
}

func (t *tUser) SetOffline(
	ctx context.Context,
	userId uuid.UUID,
) {
	key := fmt.Sprintf("user:online:%s", userId.String())

	err := t.client.Del(ctx, key).Err()
	if err != nil {
		global.Logger.Warn("Failed to delete user online key", zap.String("user_id", userId.String()), zap.Error(err))
	}
}

func (t *tUser) IsOnline(
	ctx context.Context,
	userId uuid.UUID,
) bool {
	key := fmt.Sprintf("user:online:%s", userId.String())
	val, err := t.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return val == "1"
}

func (t *tUser) ClearAllCaches(ctx context.Context) error {
	if err := t.client.FlushDB(ctx).Err(); err != nil {
		return response.NewServerFailedError(err.Error())
	}
	global.Logger.Info("Clear all caches in redis")
	return nil
}
