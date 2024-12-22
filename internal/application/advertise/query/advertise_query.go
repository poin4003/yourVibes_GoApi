package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"time"
)

type GetOneAdvertiseQuery struct {
	AdvertiseId uuid.UUID
}

type GetManyAdvertiseQuery struct {
	PostId    uuid.UUID
	Email     string
	Status    *bool
	FromDate  time.Time
	ToDate    time.Time
	FromPrice int
	ToPrice   int
	Limit     int
	Page      int
}

type GetOneAdvertiseResult struct {
	Advertise      *common.AdvertiseDetail
	ResultCode     int
	HttpStatusCode int
}

type GetManyAdvertiseResults struct {
	Advertises     []*common.AdvertiseWithBillResult
	ResultCode     int
	HttpStatusCode int
	PagingResponse *response.PagingResponse
}
