package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	reportCommand "github.com/poin4003/yourVibes_GoApi/internal/application/report/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type HandleReportRequest struct {
	ReportType consts.ReportType `json:"report_type"`
}

func ValidateHandleReportRequest(req interface{}) error {
	dto, ok := req.(*HandleReportRequest)
	if !ok {
		return fmt.Errorf("input is not Handle report request")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.ReportType, validation.In(consts.ReportTypes...)),
	)
}

func (req *HandleReportRequest) ToHandleReportCommand(
	adminId uuid.UUID,
	reportId uuid.UUID,
) (*reportCommand.HandleReportCommand, error) {
	return &reportCommand.HandleReportCommand{
		Type:     req.ReportType,
		AdminId:  adminId,
		ReportId: reportId,
	}, nil
}
