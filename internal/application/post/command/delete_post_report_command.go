package command

import "github.com/google/uuid"

type DeletePostReportCommand struct {
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
}
