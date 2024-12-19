package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
)

type HandleUserReportCommand struct {
	UserId uuid.UUID
}

type HandleUserReportCommandResult struct {
	UserReport     *common.UserReportResult
	ResultCode     int
	HttpStatusCode int
}
