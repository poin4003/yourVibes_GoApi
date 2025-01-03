package request

import (
	"github.com/google/uuid"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
)

func ToHandleUserReportCommand(
	adminId uuid.UUID,
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) (*userCommand.HandleUserReportCommand, error) {
	return &userCommand.HandleUserReportCommand{
		AdminId:        adminId,
		UserId:         userId,
		ReportedUserId: reportedUserId,
	}, nil
}
