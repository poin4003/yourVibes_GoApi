package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	admin_command "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
)

type UpdateAdminForSuperAdminRequest struct {
	AdminId *uuid.UUID `json:"admin_id"`
	Role    *bool      `json:"role"`
	Status  *bool      `json:"status"`
}

func ValidateUpdateAdminForSuperAdminRequest(req interface{}) error {
	dto, ok := req.(*UpdateAdminForSuperAdminRequest)
	if !ok {
		return fmt.Errorf("input is not UpdateAdminForSuperAdminDto")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.AdminId, validation.Required),
	)
}

func (req *UpdateAdminForSuperAdminRequest) ToUpdateAdminForSuperAdminCommand() (*admin_command.UpdateAdminForSuperAdminCommand, error) {
	return &admin_command.UpdateAdminForSuperAdminCommand{
		AdminId: req.AdminId,
		Role:    req.Role,
		Status:  req.Status,
	}, nil
}
