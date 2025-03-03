package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetOneUserQuery struct {
	UserId              uuid.UUID
	AuthenticatedUserId uuid.UUID
}

type GetManyUserQuery struct {
	Name         string
	Email        string
	PhoneNumber  string
	Birthday     time.Time
	CreatedAt    time.Time
	SortBy       string
	IsDescending bool
	Limit        int
	Page         int
}

type UserQueryResult struct {
	User       *common.UserWithoutSettingResult
	ResultCode int
}

type UserQueryListResult struct {
	Users          []*common.UserShortVerResult
	PagingResponse *response.PagingResponse
}
