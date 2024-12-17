package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	comment_command "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
)

type ReportCommentRequest struct {
	ReportCommentId uuid.UUID `json:"report_comment_id"`
	Reason          string    `json:"reason"`
}

func ValidateReportCommentRequest(req interface{}) error {
	dto, ok := req.(*ReportCommentRequest)
	if !ok {
		return fmt.Errorf("input is not ReportCommentRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.ReportCommentId, validation.Required),
		validation.Field(&dto.Reason, validation.Required, validation.Length(2, 255)),
	)
}

func (req *ReportCommentRequest) ToCreateCommentReportCommand(
	userId uuid.UUID,
) (*comment_command.CreateReportCommentCommand, error) {
	return &comment_command.CreateReportCommentCommand{
		UserId:            userId,
		ReportedCommentId: req.ReportCommentId,
		Reason:            req.Reason,
	}, nil
}
