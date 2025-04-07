package admin

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type tAdmin struct {
	client *redis.Client
}

func NewAdminCache(client *redis.Client) *tAdmin {
	return &tAdmin{client: client}
}

func (t *tAdmin) SetAdminStatus(
	ctx context.Context, adminId uuid.UUID, status bool,
) {
	key := fmt.Sprintf("admin:%s:status", adminId.String())
	val := "0"
	if status {
		val = "1"
	}
	err := t.client.Set(ctx, key, val, consts.TTL_COMMON).Err()
	if err != nil {
		global.Logger.Warn("Failed to set admin status in redis", zap.String("user_id", adminId.String()), zap.Error(err))
	}
}

func (t *tAdmin) GetAdminStatus(
	ctx context.Context, adminId uuid.UUID,
) *bool {
	key := fmt.Sprintf("admin:%s:status", adminId.String())
	val, err := t.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			global.Logger.Warn("Failed to get user status from redis", zap.Error(err))
			return nil
		}
		return nil
	}
	status := val == "1"
	return &status
}

func (t *tAdmin) DeleteAdminStatus(
	ctx context.Context, adminId uuid.UUID,
) {
	key := fmt.Sprintf("admin:%s:status", adminId.String())
	if err := t.client.Del(ctx, key).Err(); err != nil {
		global.Logger.Warn("Failed to delete user status from redis", zap.Error(err))
	}
}
