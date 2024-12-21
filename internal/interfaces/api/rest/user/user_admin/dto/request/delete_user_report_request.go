package request

import (
	"github.com/google/uuid"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
)

func ToDeleteUserReportCommand(
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) (*user_command.DeleteUserReportCommand, error) {
	return &user_command.DeleteUserReportCommand{
		UserId:         userId,
		ReportedUserId: reportedUserId,
	}, nil
}
