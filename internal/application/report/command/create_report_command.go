package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type CreateReportCommand struct {
	UserId     uuid.UUID
	ReportedId uuid.UUID
	Reason     string
	Type       consts.ReportType
}
