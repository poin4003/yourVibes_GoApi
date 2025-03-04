package common

import (
	"github.com/google/uuid"
	"time"
)

type MediaResult struct {
	ID        uint
	PostId    uuid.UUID
	MediaUrl  string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
