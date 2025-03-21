package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	reportCommand "github.com/poin4003/yourVibes_GoApi/internal/application/report/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type ActivateRequest struct {
	ReportType consts.ReportType `json:"report_type"`
}

func ValidateActivateRequest(req interface{}) error {
	dto, ok := req.(*ActivateRequest)
	if !ok {
		return fmt.Errorf("input is not Handle report request")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.ReportType, validation.In(consts.ReportTypes...)),
	)
}

func (req *ActivateRequest) ToActivateCommand(
	reportId uuid.UUID,
) (*reportCommand.ActivateCommand, error) {
	return &reportCommand.ActivateCommand{
		Type:     req.ReportType,
		ReportId: reportId,
	}, nil
}
