package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
)

type GetOnePostQuery struct {
	PostId uuid.UUID
	UserId uuid.UUID
}

type PostQueryResult struct {
	Post           *common.PostResult
	ResultCode     int
	HttpStatusCode int
}
