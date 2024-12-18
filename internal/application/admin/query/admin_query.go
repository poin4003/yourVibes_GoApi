package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"time"
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
	SortBy       string
	IsDescending bool
	Limit        int
	Page         int
}

type AdminQueryResult struct {
	Admin          *common.AdminResult
	ResultCode     int
	HttpStatusCode int
}

type AdminQueryListResult struct {
	Admins         []*common.AdminResult
	ResultCode     int
	HttpStatusCode int
	PagingResponse *response.PagingResponse
}
