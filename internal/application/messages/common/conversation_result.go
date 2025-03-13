package common

import (
	"time"

	"github.com/google/uuid"
)

type ConversationResult struct {
	ID        uuid.UUID
	Name      string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
