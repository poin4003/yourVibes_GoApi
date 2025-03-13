package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type ConversationDetail struct {
	UserId         uuid.UUID
	ConversationId uuid.UUID
	User           *User
	Conversation   *Conversation
}

func (cd *ConversationDetail) Validate() error {
	return validation.ValidateStruct(cd,
		validation.Field(&cd.UserId, validation.Required),
		validation.Field(&cd.ConversationId, validation.Required),
	)
}

func NewConversationDetail(
	UserId uuid.UUID,
	ConversationId uuid.UUID,
) (*ConversationDetail, error) {
	conversationDetail := &ConversationDetail{
		UserId:         UserId,
		ConversationId: ConversationId,
	}
	if err := conversationDetail.Validate(); err != nil {
		return nil, err
	}
	return conversationDetail, nil
}
