package command

import (
	"time"

	"github.com/google/uuid"
)

type CreateAdvertiseCommand struct {
	PostId      uuid.UUID
	StartDate   time.Time
	EndDate     time.Time
	VoucherCode *string
	RedirectUrl string
}

type CreateAdvertiseResult struct {
	PayUrl string
}
