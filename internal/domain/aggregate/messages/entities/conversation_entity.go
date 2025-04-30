package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type CreateConversation struct {
	ID                 uuid.UUID
	Name               string
	Image              string
	ConversationDetail []*ConversationDetail
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Conversation struct {
	ID             uuid.UUID
	Name           string
	Image          string
	Avatar         string
	UserID         *uuid.UUID
	FamilyName     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastMess       *string
	LastMessStatus bool
}

type ConversationUpdate struct {
	Name      *string
	Image     *string
	UpdatedAt time.Time
}

func (c *CreateConversation) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.RuneLength(1, 30)),
		validation.Field(&c.Image, is.URL),
	)
}

func (c *ConversationUpdate) ValidateConversationUpdate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.RuneLength(1, 30)),
		validation.Field(&c.Image, is.URL),
	)
}

func NewConversation(
	name string,
	userIds []uuid.UUID,
	ownerId uuid.UUID,
) (*CreateConversation, error) {
	conversationID := uuid.New()

	var conversationDetails []*ConversationDetail
	newConversationDetail, err := NewConversationDetail(ownerId, conversationID, consts.CONVERSATION_OWNER)
	if err != nil {
		return nil, err
	}
	conversationDetails = append(conversationDetails, newConversationDetail)

	for _, userId := range userIds {
		newConversationDetail, err = NewConversationDetail(userId, conversationID, consts.CONVERSATION_MEMBER)
		if err != nil {
			return nil, err
		}
		conversationDetails = append(conversationDetails, newConversationDetail)
	}

	conversation := &CreateConversation{
		ID:                 conversationID,
		Name:               name,
		Image:              consts.IMAGE_MESSAGE,
		ConversationDetail: conversationDetails,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err = conversation.Validate(); err != nil {
		return nil, err
	}

	return conversation, nil
}
