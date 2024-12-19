package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
)

type HandlePostReportCommand struct {
	UserId uuid.UUID
}

type HandlePostReportCommandResult struct {
	PostReport     *common.PostReportResult
	ResultCode     int
	HttpStatusCode int
}
