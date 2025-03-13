package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type Conversation struct {
	ID        uuid.UUID
	Name      string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Conversation) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Length(1, 30)),
		validation.Field(&c.Image, is.URL),
	)
}

func NewConversation(
	Name string,
) (*Conversation, error) {
	conversation := &Conversation{
		ID:        uuid.New(),
		Name:      Name,
		Image:     consts.IMAGE_MESSAGE,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := conversation.Validate(); err != nil {
		return nil, err
	}
	return conversation, nil

}
