package initialize

import (
	"context"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func InitRedis() *redis.Client {
	r := global.Config.Redis

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", r.Host, r.Port),
		Password: r.Password,
		DB:       r.Database,
		PoolSize: 10,
	})

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		global.Logger.Error("Redis initialization Error:", zap.Error(err))
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	global.Logger.Info("Redis initialization success")

	return rdb
}
