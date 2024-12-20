package command

import (
	"github.com/google/uuid"
)

type HandleCommentReportCommand struct {
	AdminId           uuid.UUID
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
}

type HandleCommentReportCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
