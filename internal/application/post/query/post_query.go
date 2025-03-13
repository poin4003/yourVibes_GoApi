package query

import (
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
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

type CheckPostOwnerQuery struct {
	PostId uuid.UUID
	UserId uuid.UUID
}

type GetOnePostQueryResult struct {
	Post *common.PostResultWithLiked
}

type GetManyPostQueryResult struct {
	Posts          []*common.PostResultWithLiked
	PagingResponse *response.PagingResponse
}

type CheckPostOwnerQueryResult struct {
	IsOwner bool
}
