package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetOneUserQuery struct {
	UserId              uuid.UUID
	AuthenticatedUserId uuid.UUID
}

type UserQueryResult struct {
	User           *common.UserWithoutSettingResult
	ResultCode     int
	HttpStatusCode int
}

type UserQueryListResult struct {
	Users          []*common.UserShortVerResult
	ResultCode     int
	HttpStatusCode int
	PagingResponse *pkg_response.PagingResponse
}
