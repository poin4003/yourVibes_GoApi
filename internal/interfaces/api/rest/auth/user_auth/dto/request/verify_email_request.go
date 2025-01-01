package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func ValidateVerifyEmailRequest(req interface{}) error {
	dto, ok := req.(*VerifyEmailRequest)
	if !ok {
		return fmt.Errorf("input is not VerifyEmailRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Email, validation.Required, is.Email),
	)
}
