package common

import (
	"github.com/google/uuid"
	"time"
)

type BillWithAdvertiseResult struct {
	ID          uuid.UUID
	AdvertiseId uuid.UUID
	Advertise   *AdvertiseWithoutBillResult
	Price       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      bool
}

type BillWithoutAdvertiseResult struct {
	ID          uuid.UUID
	AdvertiseId uuid.UUID
	Price       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      bool
}
