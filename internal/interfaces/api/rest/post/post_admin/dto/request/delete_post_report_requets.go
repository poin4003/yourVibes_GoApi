package request

import (
	"github.com/google/uuid"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
)

func ToDeletePostReportCommand(
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) (*post_command.DeletePostReportCommand, error) {
	return &post_command.DeletePostReportCommand{
		UserId:         userId,
		ReportedPostId: reportedPostId,
	}, nil
}
