package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ValidateChangePasswordRequest(req interface{}) error {
	dto, ok := req.(*ChangePasswordRequest)
	if !ok {
		return fmt.Errorf("input is not ChangePasswordRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.OldPassword, validation.Required, validation.RuneLength(8, 255)),
		validation.Field(&dto.NewPassword, validation.Required, validation.RuneLength(8, 255)),
	)
}

func (req *ChangePasswordRequest) ToChangePasswordCommand(
	userId uuid.UUID,
) (*userCommand.ChangePasswordCommand, error) {
	return &userCommand.ChangePasswordCommand{
		UserId:      userId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}, nil
}
