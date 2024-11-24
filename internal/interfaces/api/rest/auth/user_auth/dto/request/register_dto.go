package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"regexp"
	"time"
)

type RegisterRequest struct {
	FamilyName  string    `json:"family_name" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Email       string    `json:"email" binding:"required"`
	Password    string    `json:"password" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Birthday    time.Time `json:"birthday" binding:"required"`
	Otp         string    `json:"otp" binding:"required"`
}

func ValidateRegisterRequest(req interface{}) error {
	dto, ok := req.(*RegisterRequest)
	if !ok {
		return fmt.Errorf("input is not RegisterRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.FamilyName, validation.Required, validation.Length(2, 255)),
		validation.Field(&dto.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&dto.Email, validation.Required, is.Email),
		validation.Field(&dto.Password, validation.Required, validation.Length(8, 255)),
		validation.Field(&dto.PhoneNumber, validation.Required, validation.Length(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&dto.Birthday, validation.Required),
		validation.Field(&dto.Otp, validation.Required, validation.Length(6, 6), validation.Match((regexp.MustCompile((`^\d+$`))))),
	)
}

func (req *RegisterRequest) ToRegisterCommand() (*user_command.RegisterCommand, error) {
	return &user_command.RegisterCommand{
		FamilyName:  req.FamilyName,
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Birthday:    req.Birthday,
		Otp:         req.Otp,
	}, nil
}
