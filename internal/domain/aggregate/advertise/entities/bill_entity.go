package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Bill struct {
	ID          uuid.UUID  `validate:"omitempty,uuid4"`
	AdvertiseId uuid.UUID  `validate:"required,uuid4"`
	Advertise   *Advertise `validate:"omitempty"`
	Price       int        `validate:"required"`
	CreatedAt   time.Time  `validate:"omitempty"`
	UpdateAt    time.Time  `validate:"omitempty,gtefield=CreatedAt"`
	Status      bool       `validate:"omitempty"`
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

func (b *BillUpdate) ValidateUpdateBill() error {
	validate := validator.New()
	return validate.Struct(b)
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
