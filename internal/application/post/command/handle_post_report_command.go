package command

import (
	"github.com/google/uuid"
)

type HandlePostReportCommand struct {
	AdminId        uuid.UUID
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
}

type HandlePostReportCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
