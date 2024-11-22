package command

import (
	"github.com/google/uuid"
	"time"
)

type CreateAdvertiseCommand struct {
	PostId      uuid.UUID
	StartDate   time.Time
	EndDate     time.Time
	RedirectUrl string
}

type CreateAdvertiseResult struct {
}
