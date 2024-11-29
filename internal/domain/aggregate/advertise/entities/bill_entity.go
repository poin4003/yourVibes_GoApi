package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"time"
)

type Bill struct {
	ID          uuid.UUID
	AdvertiseId uuid.UUID
	Advertise   *Advertise
	Price       int
	CreatedAt   time.Time
	UpdateAt    time.Time
	Status      bool
}

type BillUpdate struct {
	Price  *int
	Status *bool
}

func (b *Bill) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.AdvertiseId, validation.Required),
		validation.Field(&b.Price, validation.Required),
		validation.Field(&b.UpdateAt, validation.Min(b.CreatedAt)),
	)
}

func (b *BillUpdate) ValidateUpdateBill() error {
	return validation.ValidateStruct(b)
}

func NewBill(
	AdvertiseId uuid.UUID,
	Price int,
) (*Bill, error) {
	bill := &Bill{
		ID:          uuid.New(),
		AdvertiseId: AdvertiseId,
		Price:       Price,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
		Status:      false,
	}
	if err := bill.Validate(); err != nil {
		return nil, err
	}

	return bill, nil
}
