package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/voucher/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func FromAdminModel(ad *models.Admin) *entities.Admin {
	if ad == nil {
		return nil
	}

	admin := &entities.Admin{
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
