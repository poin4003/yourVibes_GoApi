package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type FriendRequestQuery struct {
	UserId uuid.UUID
	Limit  int
	Page   int
}

type FriendRequestQueryResult struct {
	Users          []*common.UserShortVerResult
	PagingResponse *response.PagingResponse
}
