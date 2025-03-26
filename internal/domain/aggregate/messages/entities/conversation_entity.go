package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type CreateConversation struct {
	ID        uuid.UUID
	Name      string
	Image     string
	UserIds   []uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Conversation struct {
	ID        uuid.UUID
	Name      string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *CreateConversation) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Length(1, 30)),
		validation.Field(&c.Image, is.URL),
	)
}

func NewConversation(
	name string,
	userIds []uuid.UUID,
) (*CreateConversation, error) {
	conversation := &CreateConversation{
		ID:        uuid.New(),
		Name:      name,
		Image:     consts.IMAGE_MESSAGE,
		UserIds:   userIds,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := conversation.Validate(); err != nil {
		return nil, err
	}
	return conversation, nil

}
