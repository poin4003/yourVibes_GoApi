package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Advertise struct {
	ID        uuid.UUID `validate:"omitempty,uuid4"`
	PostId    uuid.UUID `validate:"required,uuid4"`
	StartDate time.Time `validate:"required"`
	EndDate   time.Time `validate:"required"`
	CreatedAt time.Time `validate:"omitempty"`
	UpdatedAt time.Time `validate:"omitempty,gtefield=CreatedAt"`
	Bill      *Bill     `validate:"omitempty"`
}

type AdvertiseUpdate struct {
	StartDate *time.Time `validate:"omitempty"`
	EndDate   *time.Time `validate:"omitempty"`
	UpdatedAt *time.Time `validate:"omitempty,gtefield=CreatedAt"`
}

func (a *Advertise) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}

func (a *AdvertiseUpdate) ValidateAdvertiseUpdate() error {
	validate := validator.New()
	return validate.Struct(a)
}

func NewAdvertise(
	PostId uuid.UUID,
	StartDate time.Time,
	EndDate time.Time,
) (*Advertise, error) {
	advertise := &Advertise{
		ID:        uuid.New(),
		PostId:    PostId,
		StartDate: StartDate,
		EndDate:   EndDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := advertise.Validate(); err != nil {
		return nil, err
	}

	return advertise, nil
}
