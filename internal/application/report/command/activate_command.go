package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type ActivateCommand struct {
	ReportId uuid.UUID
	Type     consts.ReportType
}
