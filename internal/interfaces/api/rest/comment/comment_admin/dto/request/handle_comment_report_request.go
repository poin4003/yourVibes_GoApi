package request

import (
	"github.com/google/uuid"
	commentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
)

func ToHandleCommentReportCommand(
	adminId uuid.UUID,
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
) (*commentCommand.HandleCommentReportCommand, error) {
	return &commentCommand.HandleCommentReportCommand{
		AdminId:           adminId,
		UserId:            userId,
		ReportedCommentId: reportedCommentId,
	}, nil
}
