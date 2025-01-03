package request

import (
	"github.com/google/uuid"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
)

func ToDeleteUserReportCommand(
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) (*userCommand.DeleteUserReportCommand, error) {
	return &userCommand.DeleteUserReportCommand{
		UserId:         userId,
		ReportedUserId: reportedUserId,
	}, nil
}
