package query

import (
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
)

type GetOneAdvertiseQuery struct {
	AdvertiseId uuid.UUID
}

type GetManyAdvertiseQuery struct {
	PostId       uuid.UUID
	UserEmail    string
	Status       *bool
	FromDate     time.Time
	ToDate       time.Time
	FromPrice    int
	ToPrice      int
	SortBy       string
	IsDescending bool
	Limit        int
	Page         int
}

type GetManyAdvertiseByUserId struct {
	UserId uuid.UUID
	Limit  int
	Page   int
}

type GetOneAdvertiseResult struct {
	Advertise *common.AdvertiseDetailResult
}

type GetManyAdvertiseResults struct {
	Advertises     []*common.AdvertiseWithBillResult
	PagingResponse *response.PagingResponse
}

type GetManyAdvertiseResultsByUserId struct {
	Advertises     []*common.ShortAdvertiseResult
	PagingResponse *response.PagingResponse
}
