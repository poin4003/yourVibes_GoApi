package common

import (
	"time"

	"github.com/google/uuid"
)

type AdvertiseWithBillResult struct {
	ID           uuid.UUID
	PostId       uuid.UUID
	UserEmail    string
	StartDate    time.Time
	EndDate      time.Time
	DayRemaining int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Bill         *BillWithoutAdvertiseResult
}

type AdvertiseWithoutBillResult struct {
	ID           uuid.UUID
	PostId       uuid.UUID
	StartDate    time.Time
	EndDate      time.Time
	DayRemaining int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AdvertiseDetailResult struct {
	ID           uuid.UUID
	PostId       uuid.UUID
	Post         *PostForAdvertiseResult
	StartDate    time.Time
	EndDate      time.Time
	DayRemaining int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Bill         *BillWithoutAdvertiseResult
}

type AdvertiseForStatisticResult struct {
	ID              uuid.UUID
	PostId          uuid.UUID
	StartDate       time.Time
	EndDate         time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DayRemaining    int
	Bill            *BillWithoutAdvertiseResult
	Post            *PostForAdvertiseResult
	TotalReach      int
	TotalClicks     int
	TotalImpression int
	Statistics      []*StatisticResult
}
