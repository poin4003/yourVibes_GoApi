package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type GetOtpForgotUserPasswordRequest struct {
	Email string `json:"email"`
}

func ValidateGetOtpForgotUserPasswordRequest(req interface{}) error {
	dto, ok := req.(*GetOtpForgotUserPasswordRequest)
	if !ok {
		return fmt.Errorf("input is not GetOtpForgotUserPasswordRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Email, validation.Required, is.Email),
	)
}
