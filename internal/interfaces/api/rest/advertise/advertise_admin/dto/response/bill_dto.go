package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"time"
)

type BillWithoutAdvertiseDto struct {
	ID          uuid.UUID `json:"id"`
	AdvertiseId uuid.UUID `json:"advertise_id"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      bool      `json:"status"`
}

func ToBillWithoutAdvertiseDto(
	billResult common.BillWithoutAdvertiseResult,
) *BillWithoutAdvertiseDto {
	billDto := &BillWithoutAdvertiseDto{
		ID:          billResult.ID,
		AdvertiseId: billResult.AdvertiseId,
		Price:       billResult.Price,
		CreatedAt:   billResult.CreatedAt,
		UpdatedAt:   billResult.UpdatedAt,
		Status:      billResult.Status,
	}
	return billDto
}
