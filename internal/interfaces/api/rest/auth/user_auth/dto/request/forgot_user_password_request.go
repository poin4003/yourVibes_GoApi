package request

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
)

type ForgotUserPasswordRequest struct {
	Email       string `json:"email"`
	Otp         string `json:"otp"`
	NewPassword string `json:"new_password"`
}

func ValidateForgotUserPasswordRequest(req interface{}) error {
	dto, ok := req.(*ForgotUserPasswordRequest)
	if !ok {
		return fmt.Errorf("input is not ForgotUserPasswordRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Email, validation.Required, is.Email),
		validation.Field(&dto.Otp, validation.Required, validation.RuneLength(6, 6), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&dto.NewPassword, validation.Required, validation.RuneLength(8, 255)),
	)
}

func (req *ForgotUserPasswordRequest) ToForgotUserPasswordCommand() (*userCommand.ForgotUserPasswordCommand, error) {
	return &userCommand.ForgotUserPasswordCommand{
		Email:       req.Email,
		Otp:         req.Otp,
		NewPassword: req.NewPassword,
	}, nil
}
