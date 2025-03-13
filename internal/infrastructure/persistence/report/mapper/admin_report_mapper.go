package mapper

import (
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func FromAdminModel(ad *models.Admin) *reportEntity.Admin {
	if ad == nil {
		return nil
	}

	admin := &reportEntity.Admin{
		FamilyName:  ad.FamilyName,
		Name:        ad.Name,
		Email:       ad.Email,
		PhoneNumber: ad.PhoneNumber,
		IdentityId:  ad.IdentityId,
		Birthday:    ad.Birthday,
		Status:      ad.Status,
		Role:        ad.Role,
		CreatedAt:   ad.CreatedAt,
		UpdatedAt:   ad.UpdatedAt,
	}
	admin.ID = ad.ID

	return admin
}
