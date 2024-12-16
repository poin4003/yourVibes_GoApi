package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/common"
	"time"
)

type AdminDto struct {
	ID          uuid.UUID
	FamilyName  string
	Name        string
	Email       string
	PhoneNumber string
	IdentityId  string
	Birthday    time.Time
	Status      bool
	Role        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ToAdminDto(
	adminResult *common.AdminResult,
) *AdminDto {
	return &AdminDto{
		ID:          adminResult.ID,
		FamilyName:  adminResult.FamilyName,
		Name:        adminResult.Name,
		Email:       adminResult.Email,
		PhoneNumber: adminResult.PhoneNumber,
		IdentityId:  adminResult.IdentityId,
		Birthday:    adminResult.Birthday,
		Status:      adminResult.Status,
		Role:        adminResult.Role,
		CreatedAt:   adminResult.CreatedAt,
		UpdatedAt:   adminResult.UpdatedAt,
	}
}
