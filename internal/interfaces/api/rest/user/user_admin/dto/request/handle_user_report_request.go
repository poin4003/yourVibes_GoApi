package request

import (
	"github.com/google/uuid"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
)

func ToHandleUserReportCommand(
	adminId uuid.UUID,
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) (*user_command.HandleUserReportCommand, error) {
	return &user_command.HandleUserReportCommand{
		AdminId:        adminId,
		UserId:         userId,
		ReportedUserId: reportedUserId,
	}, nil
}
