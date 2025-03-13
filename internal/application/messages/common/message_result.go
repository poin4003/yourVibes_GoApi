package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageResult struct {
	ID             uuid.UUID
	UserId         uuid.UUID
	User           *UserResult
	ConversationId uuid.UUID
	ParentId       *uuid.UUID
	Content        *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}
