package entities

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a message entity.
type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
	IsRead         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// Conversation represents a conversation entity.
type Conversation struct {
	ID        uuid.UUID
	Name      *string // Optional name for the conversation (e.g., group chat name)
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// MessageUser represents a user in a conversation (many-to-many relationship).
type MessageUser struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID
	JoinedAt       time.Time
	LeftAt         *time.Time
}

// MessageLike represents like for message
type MessageLike struct {
	MessageID uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}
