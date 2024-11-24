package common

import (
	"github.com/google/uuid"
	"time"
)

type BillWithAdvertiseResult struct {
	ID          uuid.UUID
	AdvertiseId uuid.UUID
	Advertise   *AdvertiseWithoutBillResult
	Price       float64
	Vat         float64
	CreatedAt   time.Time
	UpdateAt    time.Time
	Status      bool
}

type BillWithoutAdvertiseResult struct {
	ID          uuid.UUID
	AdvertiseId uuid.UUID
	Price       float64
	Vat         float64
	CreatedAt   time.Time
	UpdateAt    time.Time
	Status      bool
}
