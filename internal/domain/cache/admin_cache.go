package cache

import (
	"context"
	"github.com/google/uuid"
)

type (
	IAdminCache interface {
		SetAdminStatus(ctx context.Context, adminId uuid.UUID, status bool)
		GetAdminStatus(ctx context.Context, adminId uuid.UUID) *bool
		DeleteAdminStatus(ctx context.Context, adminId uuid.UUID)
	}
)

var (
	localAdminCache IAdminCache
)

func InitAdminCache(i IAdminCache) {
	localAdminCache = i
}
