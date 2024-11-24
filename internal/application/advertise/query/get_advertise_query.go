package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
)

type GetOneAdvertiseQuery struct {
	PostId uuid.UUID
}

type GetOneAdvertiseResult struct {
	Advertise      *common.AdvertiseWithBillResult
	ResultCode     int
	HttpStatusCode int
}
