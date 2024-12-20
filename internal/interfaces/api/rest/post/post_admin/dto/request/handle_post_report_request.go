package request

import (
	"github.com/google/uuid"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
)

func ToHandlePostReportCommand(
	adminId uuid.UUID,
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) (*post_command.HandlePostReportCommand, error) {
	return &post_command.HandlePostReportCommand{
		AdminId:        adminId,
		UserId:         userId,
		ReportedPostId: reportedPostId,
	}, nil
}
