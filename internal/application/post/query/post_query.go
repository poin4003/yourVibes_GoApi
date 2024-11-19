package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"time"
)

type GetOnePostQuery struct {
	PostId              uuid.UUID
	AuthenticatedUserId uuid.UUID
}

type GetManyPostQuery struct {
	AuthenticatedUserId uuid.UUID
	UserID              uuid.UUID
	Content             string
	Location            string
	IsAdvertisement     bool
	CreatedAt           time.Time
	SortBy              string
	IsDescending        bool
	Limit               int
	Page                int
}

type GetOnePostQueryResult struct {
	Post           *common.PostResultWithLiked
	ResultCode     int
	HttpStatusCode int
}

type GetManyPostQueryResult struct {
	Posts          []*common.PostResultWithLiked
	ResultCode     int
	HttpStatusCode int
	PagingResponse *response.PagingResponse
}
