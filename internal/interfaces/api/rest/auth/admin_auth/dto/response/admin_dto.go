package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/common"
	"time"
)

type AdminDto struct {
	ID          uuid.UUID `json:"id"`
	FamilyName  string    `json:"family_name"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	IdentityId  string    `json:"identity_id"`
	Birthday    time.Time `json:"birthday"`
	Status      bool      `json:"status"`
	Role        bool      `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
