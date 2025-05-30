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
		SetOnline(ctx context.Context, userId uuid.UUID)
		SetOffline(ctx context.Context, userId uuid.UUID)
		IsOnline(ctx context.Context, userId uuid.UUID) bool
		ClearAllCaches(ctx context.Context) error
	}
	IUserAuthCache interface {
		SetOtp(ctx context.Context, userKey, otp string, ttl time.Duration) error
		GetOtp(ctx context.Context, userKey string) (*string, error)
	}
)
