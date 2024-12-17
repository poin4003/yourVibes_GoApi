package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
)

type ReportUserRequest struct {
	ReportedUserId uuid.UUID `json:"reported_user_id"`
	Reason         string    `json:"reason"`
}

func ValidateReportUserRequest(req interface{}) error {
	dto, ok := req.(*ReportUserRequest)
	if !ok {
		return fmt.Errorf("input is not ReportUserRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.ReportedUserId, validation.Required),
		validation.Field(&dto.Reason, validation.Required, validation.Length(2, 255)),
	)
}

func (req *ReportUserRequest) ToCreateUserReportCommand(
	userId uuid.UUID,
) (*user_command.CreateReportUserCommand, error) {
	return &user_command.CreateReportUserCommand{
		UserId:         userId,
		ReportedUserId: req.ReportedUserId,
		Reason:         req.Reason,
	}, nil
}
