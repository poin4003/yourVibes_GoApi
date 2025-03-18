package common

import (
	"time"

	"github.com/google/uuid"
)

type MessageResult struct {
	ID             uuid.UUID
	UserId         uuid.UUID
	User           *UserResult
	ConversationId uuid.UUID
	ParentId       *uuid.UUID
	ParentContent  *string
	Content        *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
