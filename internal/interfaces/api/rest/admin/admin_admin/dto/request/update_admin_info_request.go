package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	admin_command "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"regexp"
	"time"
)

type UpdateAdminInfoRequest struct {
	FamilyName  string    `json:"family_name"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	IdentityId  string    `json:"identity_id"`
	Birthday    time.Time `json:"birthday"`
}

func ValidateUpdateAdminInfoRequest(req interface{}) error {
	dto, ok := req.(*UpdateAdminInfoRequest)
	if !ok {
		return fmt.Errorf("input is not UpdateAdminRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.FamilyName, validation.Length(2, 255)),
		validation.Field(&dto.Name, validation.Length(2, 255)),
		validation.Field(&dto.PhoneNumber, validation.Length(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&dto.IdentityId, validation.Length(10, 15), validation.Match((regexp.MustCompile((`^\d+$`))))),
	)
}

func (req *UpdateAdminInfoRequest) ToUpdateAdminInfoCommand(
	adminId uuid.UUID,
) (*admin_command.UpdateAdminInfoCommand, error) {
	return &admin_command.UpdateAdminInfoCommand{
		AdminID:     &adminId,
		FamilyName:  &req.FamilyName,
		Name:        &req.Name,
		PhoneNumber: &req.PhoneNumber,
		IdentityId:  &req.IdentityId,
		Birthday:    &req.Birthday,
	}, nil
}
