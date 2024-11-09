package request

import (
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"time"
)

type RegisterRequest struct {
	FamilyName  string    `json:"family_name" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Password    string    `json:"password" binding:"required,min=8"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Birthday    time.Time `json:"birthday" binding:"required"`
	Otp         string    `json:"otp" binding:"required"`
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
