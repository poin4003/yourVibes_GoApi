package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/common"
	adminEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	adminValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/validator"
)

func NewAdminResultFromValidateEntity(
	admin *adminValidator.ValidateAdmin,
) *common.AdminResult {
	return NewAdminResult(&admin.Admin)
}

func NewAdminResult(
	admin *adminEntity.Admin,
) *common.AdminResult {
	if admin == nil {
		return nil
	}

	return &common.AdminResult{
		ID:          admin.ID,
		FamilyName:  admin.FamilyName,
		Name:        admin.Name,
		Email:       admin.Email,
		PhoneNumber: admin.PhoneNumber,
		IdentityId:  admin.IdentityId,
		Birthday:    admin.Birthday,
		Status:      admin.Status,
		Role:        admin.Role,
		CreatedAt:   admin.CreatedAt,
		UpdatedAt:   admin.UpdatedAt,
	}
}
