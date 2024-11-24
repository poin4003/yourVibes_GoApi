package common

import (
	"github.com/google/uuid"
	"time"
)

type AdvertiseWithBillResult struct {
	ID        uuid.UUID
	PostId    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Bill      *BillWithoutAdvertiseResult
}

type AdvertiseWithoutBillResult struct {
	ID        uuid.UUID
	PostId    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
