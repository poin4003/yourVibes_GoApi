package request

import (
	"github.com/google/uuid"
	commentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
)

func ToDeleteCommentReportCommand(
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
) (*commentCommand.DeleteCommentReportCommand, error) {
	return &commentCommand.DeleteCommentReportCommand{
		UserId:            userId,
		ReportedCommentId: reportedCommentId,
	}, nil
}
