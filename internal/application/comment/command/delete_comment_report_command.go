package command

import "github.com/google/uuid"

type DeleteCommentReportCommand struct {
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
}
