package entities

import (
	"github.com/google/uuid"
	"time"
)

type Media struct {
	ID        uint
	PostId    uuid.UUID
	MediaUrl  string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
