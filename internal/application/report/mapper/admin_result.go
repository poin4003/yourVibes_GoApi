package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
)

func NewAdminResult(
	admin *reportEntity.Admin,
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
