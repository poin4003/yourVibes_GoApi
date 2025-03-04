package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

// ValidatedMessage is a struct that holds a message to be validated.
type ValidatedMessage struct {
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
}

// Validate validates the ValidatedMessage struct.
func (m *ValidatedMessage) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.ConversationID, validation.Required),
		validation.Field(&m.SenderID, validation.Required),
		validation.Field(&m.Content, validation.Required, validation.Length(1, 10000)), // Assuming a max content length of 10000 characters
	)
}
