package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
)

type ReportPostRequest struct {
	ReportPostId uuid.UUID `json:"report_post_id"`
	Reason       string    `json:"reason"`
}

func ValidateReportPostRequest(req interface{}) error {
	dto, ok := req.(*ReportPostRequest)
	if !ok {
		return fmt.Errorf("input is not ReportPostRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.ReportPostId, validation.Required),
		validation.Field(&dto.Reason, validation.Required, validation.Length(2, 255)),
	)
}

func (req *ReportPostRequest) ToCreatePostReportCommand(
	userId uuid.UUID,
) (*post_command.CreateReportPostCommand, error) {
	return &post_command.CreateReportPostCommand{
		UserId:         userId,
		ReportedPostId: req.ReportPostId,
		Reason:         req.Reason,
	}, nil
}
