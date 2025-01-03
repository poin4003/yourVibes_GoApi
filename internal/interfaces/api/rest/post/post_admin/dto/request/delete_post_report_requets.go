package request

import (
	"github.com/google/uuid"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
)

func ToDeletePostReportCommand(
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) (*postCommand.DeletePostReportCommand, error) {
	return &postCommand.DeletePostReportCommand{
		UserId:         userId,
		ReportedPostId: reportedPostId,
	}, nil
}
