package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
)

type CreateReportCommentCommand struct {
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
	Reason            string
}

type CreateReportCommentCommandResult struct {
	CommandReport  *common.UserReportResult
	ResultCode     int
	HttpStatusCode int
}