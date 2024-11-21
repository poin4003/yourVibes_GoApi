package entities

import (
	"github.com/google/uuid"
	"time"
)

type Advertise struct {
	ID        uuid.UUID
	PostId    uuid.UUID
	TotalView int
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
