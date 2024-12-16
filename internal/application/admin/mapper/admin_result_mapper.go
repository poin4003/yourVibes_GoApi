package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	admin_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	admin_validator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/validator"
)

func NewAdminValidateEntity(
	admin *admin_validator.ValidateAdmin,
) *common.AdminResult {
	return NewAdminEntity(&admin.Admin)
}

func NewAdminEntity(
	admin *admin_entity.Admin,
) *common.AdminResult {
	if admin == nil {
		return nil
	}

	return &common.AdminResult{
		ID:          admin.ID,
		FamilyName:  admin.FamilyName,
		Name:        admin.Name,
		Email:       admin.Email,
		Password:    admin.Password,
		PhoneNumber: admin.PhoneNumber,
		IdentityId:  admin.IdentityId,
		Birthday:    admin.Birthday,
		Status:      admin.Status,
		Role:        admin.Role,
		CreatedAt:   admin.CreatedAt,
		UpdatedAt:   admin.UpdatedAt,
	}
}
