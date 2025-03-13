package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type HandleReportCommand struct {
	ReportId uuid.UUID
	AdminId  uuid.UUID
	Type     consts.ReportType
}
