package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
)

type AdminDto struct {
	ID         uuid.UUID `json:"id"`
	FamilyName string    `json:"family_name"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Status     bool      `json:"status"`
}

func ToAdminDto(
	adminResult *common.AdminResult,
) *AdminDto {
	if adminResult == nil {
		return nil
	}

	return &AdminDto{
		ID:         adminResult.ID,
		FamilyName: adminResult.FamilyName,
		Name:       adminResult.Name,
		Email:      adminResult.Email,
		Status:     adminResult.Status,
	}
}
