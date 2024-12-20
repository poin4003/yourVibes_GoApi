package request

import (
	"github.com/google/uuid"
	comment_command "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
)

func ToHandleCommentReportCommand(
	adminId uuid.UUID,
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
) (*comment_command.HandleCommentReportCommand, error) {
	return &comment_command.HandleCommentReportCommand{
		AdminId:           adminId,
		UserId:            userId,
		ReportedCommentId: reportedCommentId,
	}, nil
}
