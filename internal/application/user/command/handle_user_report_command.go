package command

import (
	"github.com/google/uuid"
)

type HandleUserReportCommand struct {
	AdminId        uuid.UUID
	UserId         uuid.UUID
	ReportedUserId uuid.UUID
}

type HandleUserReportCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}