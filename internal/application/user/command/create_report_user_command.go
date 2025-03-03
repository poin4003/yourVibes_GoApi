package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
)

type CreateReportUserCommand struct {
	UserId         uuid.UUID
	ReportedUserId uuid.UUID
	Reason         string
}

type CreateReportUserCommandResult struct {
	UserReport *common.UserReportResult
}
