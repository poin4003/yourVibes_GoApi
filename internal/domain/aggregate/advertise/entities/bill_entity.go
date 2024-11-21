package entities

import (
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
	Status      bool
}
