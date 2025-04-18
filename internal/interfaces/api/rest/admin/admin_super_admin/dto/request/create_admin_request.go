package request

import (
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	adminCommand "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
)

type CreateAdminRequest struct {
	FamilyName  string    `json:"family_name"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone_number"`
	IdentityId  string    `json:"identity_id"`
	Birthday    time.Time `json:"birthday"`
	Role        *bool     `json:"role,omitempty"`
}

func ValidateCreateAdminRequest(req interface{}) error {
	dto, ok := req.(*CreateAdminRequest)
	if !ok {
		return fmt.Errorf("input is not CreateAdminRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.FamilyName, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&dto.Name, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&dto.Email, validation.Required, is.Email),
		validation.Field(&dto.Password, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&dto.PhoneNumber, validation.Required, validation.RuneLength(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&dto.IdentityId, validation.Required, validation.RuneLength(10, 15), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&dto.Birthday, validation.Required),
	)
}

func (req *CreateAdminRequest) ToCreateAdminCommand() (*adminCommand.CreateAdminCommand, error) {
	return &adminCommand.CreateAdminCommand{
		FamilyName:  req.FamilyName,
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		IdentityId:  req.IdentityId,
		Birthday:    req.Birthday,
		Role:        *req.Role,
	}, nil
}
