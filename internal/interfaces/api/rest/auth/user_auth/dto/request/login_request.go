package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func ValidateLoginRequest(req interface{}) error {
	dto, ok := req.(*LoginRequest)
	if !ok {
		return fmt.Errorf("input is not LoginRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Email, validation.Required, is.Email),
		validation.Field(&dto.Password, validation.Required, validation.RuneLength(8, 255)),
	)
}

func (req *LoginRequest) ToLoginCommand() (*userCommand.LoginCommand, error) {
	return &userCommand.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}
