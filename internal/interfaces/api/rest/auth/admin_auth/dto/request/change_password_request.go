package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	adminCommand "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
)

type ChangeAdminPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ValidateChangePasswordRequest(req interface{}) error {
	dto, ok := req.(*ChangeAdminPasswordRequest)
	if !ok {
		return fmt.Errorf("input is not ChangeAdminPasswordRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.OldPassword, validation.Required, validation.Length(8, 255)),
		validation.Field(&dto.NewPassword, validation.Required, validation.Length(8, 255)),
	)
}

func (req *ChangeAdminPasswordRequest) ToChangeAdminPasswordCommand(
	adminId uuid.UUID,
) (*adminCommand.ChangeAdminPasswordCommand, error) {
	return &adminCommand.ChangeAdminPasswordCommand{
		AdminId:     adminId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}, nil
}
