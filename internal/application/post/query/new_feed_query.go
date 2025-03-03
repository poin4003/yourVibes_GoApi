package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetNewFeedQuery struct {
	UserId uuid.UUID
	Limit  int
	Page   int
}

type GetNewFeedResult struct {
	Posts          []*common.PostResultWithLiked
	PagingResponse *response.PagingResponse
}
