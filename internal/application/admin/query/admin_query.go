package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetOneAdminQuery struct {
	AdminId uuid.UUID
}

type GetManyAdminQuery struct {
	Name         string
	Email        string
	PhoneNumber  string
	IdentityId   string
	Birthday     time.Time
	CreatedAt    time.Time
	Status       *bool
	Role         *bool
	SortBy       string
	IsDescending bool
	Limit        int
	Page         int
}

type AdminQueryResult struct {
	Admin *common.AdminResult
}

type AdminQueryListResult struct {
	Admins         []*common.AdminResult
	PagingResponse *response.PagingResponse
}
