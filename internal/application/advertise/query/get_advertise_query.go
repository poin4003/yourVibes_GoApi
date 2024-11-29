package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetOneAdvertiseQuery struct {
	PostId uuid.UUID
}

type GetManyAdvertiseQuery struct {
	PostId uuid.UUID
	Limit  int
	Page   int
}

type GetOneAdvertiseResult struct {
	Advertise      *common.AdvertiseWithBillResult
	ResultCode     int
	HttpStatusCode int
}

type GetManyAdvertiseResults struct {
	Advertises     []*common.AdvertiseWithBillResult
	ResultCode     int
	HttpStatusCode int
	PagingResponse *response.PagingResponse
}
