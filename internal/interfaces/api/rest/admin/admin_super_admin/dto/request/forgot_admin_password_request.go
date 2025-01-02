package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	admin_command "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
)

type ForgotAdminPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

func ValidateForgotAdminPasswordRequest(req interface{}) error {
	dto, ok := req.(*ForgotAdminPasswordRequest)
	if !ok {
		return fmt.Errorf("input is not ForgotAdminPasswordRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Email, validation.Required, is.Email),
		validation.Field(&dto.NewPassword, validation.Required, validation.Length(8, 255)),
	)
}

func (req *ForgotAdminPasswordRequest) ToForgotAdminPasswordCommand() (*admin_command.ForgotAdminPasswordCommand, error) {
	return &admin_command.ForgotAdminPasswordCommand{
		Email:       req.Email,
		NewPassword: req.NewPassword,
	}, nil
}
