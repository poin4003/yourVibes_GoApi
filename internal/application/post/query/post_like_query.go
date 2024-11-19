package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetPostLikeQuery struct {
	PostId uuid.UUID
	Limit  int
	Page   int
}

type GetPostLikeQueryResult struct {
	Users          []*common.UserResult
	ResultCode     int
	HttpStatusCode int
	PagingResponse *response.PagingResponse
}
