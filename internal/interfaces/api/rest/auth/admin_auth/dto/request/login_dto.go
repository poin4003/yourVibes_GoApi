package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	admin_command "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
)

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func ValidateLoginRequest(req interface{}) error {
	dto, ok := req.(*AdminLoginRequest)
	if !ok {
		return fmt.Errorf("input is not LoginRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Email, validation.Required, is.Email),
		validation.Field(&dto.Password, validation.Required, validation.Length(8, 255)),
	)
}

func (req *AdminLoginRequest) ToLoginCommand() (*admin_command.LoginCommand, error) {
	return &admin_command.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}
