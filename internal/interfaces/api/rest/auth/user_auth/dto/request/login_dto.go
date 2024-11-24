package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
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
		validation.Field(&dto.Password, validation.Required, validation.Length(8, 255)),
	)
}

func (req *LoginRequest) ToLoginCommand() (*user_command.LoginCommand, error) {
	return &user_command.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}
