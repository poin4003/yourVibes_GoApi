package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Bill struct {
	ID          uuid.UUID
	AdvertiseId uuid.UUID
	Advertise   *Advertise
	Price       float64
	Vat         float64
	CreatedAt   time.Time
	UpdateAt    time.Time
	Status      bool
}

type BillUpdate struct {
	Price  *float64
	Vat    *float64
	Status *bool
}

func (b *Bill) Validate() error {
	validate := validator.New()
	return validate.Struct(b)
}

func (b *Bill) ValidateUpdateBill() error {
	validate := validator.New()
	return validate.Struct(b)
}

func NewBill(
	AdvertiseId uuid.UUID,
	Price float64,
) (*Bill, error) {
	bill := &Bill{
		ID:          uuid.New(),
		AdvertiseId: AdvertiseId,
		Price:       Price,
		Vat:         0.1,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
		Status:      false,
	}
	if err := bill.Validate(); err != nil {
		return nil, err
	}

	return bill, nil
}
