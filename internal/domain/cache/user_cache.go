package cache

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type (
	IUserCache interface {
		SetUserStatus(ctx context.Context, userId uuid.UUID, status bool)
		GetUserStatus(ctx context.Context, userId uuid.UUID) *bool
		DeleteUserStatus(ctx context.Context, userId uuid.UUID)
	}
	IUserAuthCache interface {
		SetOtp(ctx context.Context, userKey, otp string, ttl time.Duration) error
		GetOtp(ctx context.Context, userKey string) (*string, error)
	}
)

var (
	localUserAuthCache IUserAuthCache
	localUserCache     IUserCache
)

func UserAuthCache() IUserAuthCache {
	if localUserAuthCache == nil {
		panic("repository_implement localUserAuth not found for interface IUserAuthCache")
	}

	return localUserAuthCache
}

func UserCache() IUserCache {
	if localUserCache == nil {
		panic("repository_implement localUserAuth not found for interface IUserAuthCache")
	}

	return localUserCache
}

func InitUserAuthCache(i IUserAuthCache) {
	localUserAuthCache = i
}

func InitUserCache(i IUserCache) {
	localUserCache = i
}
