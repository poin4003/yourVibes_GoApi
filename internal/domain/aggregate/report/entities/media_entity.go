package entities

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID        uint
	PostId    uuid.UUID
	MediaUrl  string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
