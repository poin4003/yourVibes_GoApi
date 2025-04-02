package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID
	UserId         uuid.UUID
	User           *User
	ConversationId uuid.UUID
	ParentId       *uuid.UUID
	ParentContent  *string
	Content        *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (m *Message) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Content, validation.Required, validation.RuneLength(1, 1000)),
		validation.Field(&m.ConversationId, validation.Required),
		validation.Field(&m.UserId, validation.Required),
	)
}

func NewMessage(
	userId uuid.UUID,
	conversationId uuid.UUID,
	parentId *uuid.UUID,
	content *string,
) (*Message, error) {
	message := &Message{
		ID:             uuid.New(),
		UserId:         userId,
		ConversationId: conversationId,
		ParentId:       parentId,
		Content:        content,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}
