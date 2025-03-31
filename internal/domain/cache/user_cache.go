package cache

import (
	"context"
	"time"
)

type (
	IUserAuthCache interface {
		SetOtp(ctx context.Context, userKey, otp string, ttl time.Duration) error
		GetOtp(ctx context.Context, userKey string) (*string, error)
	}
)

var (
	localUserAuthCache IUserAuthCache
)

func UserAuthCache() IUserAuthCache {
	if localUserAuthCache == nil {
		panic("repository_implement localUserAuth not found for interface IUserAuthCache")
	}

	return localUserAuthCache
}

func InitUserAuthCache(i IUserAuthCache) {
	localUserAuthCache = i
}
