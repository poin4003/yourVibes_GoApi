package request

import (
	"github.com/google/uuid"
	comment_command "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
)

func ToDeleteCommentReportCommand(
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
) (*comment_command.DeleteCommentReportCommand, error) {
	return &comment_command.DeleteCommentReportCommand{
		UserId:            userId,
		ReportedCommentId: reportedCommentId,
	}, nil
}
