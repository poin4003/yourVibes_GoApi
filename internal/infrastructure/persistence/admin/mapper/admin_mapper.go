package mapper

import (
	adminEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToAdminModel(admin *adminEntity.Admin) *models.Admin {
	ad := &models.Admin{
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
	ad.ID = admin.ID

	return ad
}

func FromAdminModel(ad *models.Admin) *adminEntity.Admin {
	admin := &adminEntity.Admin{
		FamilyName:  ad.FamilyName,
		Name:        ad.Name,
		Email:       ad.Email,
		Password:    ad.Password,
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
