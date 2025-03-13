package command

import "github.com/google/uuid"

type DeleteReportCommand struct {
	ReportId uuid.UUID
}
