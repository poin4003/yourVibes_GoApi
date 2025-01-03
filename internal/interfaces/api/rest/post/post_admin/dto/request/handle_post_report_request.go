package request

import (
	"github.com/google/uuid"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
)

func ToHandlePostReportCommand(
	adminId uuid.UUID,
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) (*postCommand.HandlePostReportCommand, error) {
	return &postCommand.HandlePostReportCommand{
		AdminId:        adminId,
		UserId:         userId,
		ReportedPostId: reportedPostId,
	}, nil
}
