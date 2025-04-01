package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/report/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type ReportRequest struct {
	ReportedId uuid.UUID         `json:"reported_id"`
	Reason     string            `json:"reason"`
	Type       consts.ReportType `json:"type"`
}

func ValidateReportRequest(req interface{}) error {
	dto, ok := req.(*ReportRequest)
	if !ok {
		return fmt.Errorf("input is not ReportRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.ReportedId, validation.Required),
		validation.Field(&dto.Reason, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&dto.Type, validation.In(consts.ReportTypes...)),
	)
}

func (req *ReportRequest) ToCreateReportCommand(
	userId uuid.UUID,
) (*userCommand.CreateReportCommand, error) {
	return &userCommand.CreateReportCommand{
		UserId:     userId,
		ReportedId: req.ReportedId,
		Reason:     req.Reason,
		Type:       req.Type,
	}, nil
}
