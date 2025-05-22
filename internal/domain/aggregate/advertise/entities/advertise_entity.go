package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type Advertise struct {
	ID        uuid.UUID
	PostId    uuid.UUID
	Post      *PostForAdvertise
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Bill      *Bill
}

type ShortAdvertise struct {
	Post      *ShortPostForAdvertise
	StartDate time.Time
	EndDate   time.Time
	BillPrice int
}

type AdvertiseForStatistic struct {
	ID               uuid.UUID
	PostId           uuid.UUID
	StartDate        time.Time
	EndDate          time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Bill             *Bill
	PostForAdvertise *PostForAdvertise
	TotalReach       int
	TotalClicks      int
	TotalImpression  int
	Statistics       []*StatisticEntity
}

type AdvertiseUpdate struct {
	StartDate *time.Time
	EndDate   *time.Time
	UpdatedAt *time.Time
}

func (a *Advertise) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.PostId, validation.Required),
		validation.Field(&a.StartDate, validation.Required),
		validation.Field(&a.EndDate, validation.Required),
		validation.Field(&a.UpdatedAt, validation.Min(a.CreatedAt)),
	)
}

func (a *AdvertiseUpdate) ValidateAdvertiseUpdate() error {
	return validation.ValidateStruct(a)
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
