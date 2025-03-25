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
	LastMessStatus bool
}

type ConversationDetailUpdate struct {
	LastMessStatus *bool
}

func (cd *ConversationDetail) Validate() error {
	return validation.ValidateStruct(cd,
		validation.Field(&cd.UserId, validation.Required),
		validation.Field(&cd.ConversationId, validation.Required),
		validation.Field(&cd.LastMessStatus, validation.Required),
	)
}

func NewConversationDetail(
	UserId uuid.UUID,
	ConversationId uuid.UUID,
) (*ConversationDetail, error) {
	conversationDetail := &ConversationDetail{
		UserId:         UserId,
		ConversationId: ConversationId,
		LastMessStatus: true,
	}
	if err := conversationDetail.Validate(); err != nil {
		return nil, err
	}
	return conversationDetail, nil
}

func NewConversationDetailUpdate(
	updateData *ConversationDetailUpdate,
) (*ConversationDetailUpdate, error) {
	conversationDetailUpdate := &ConversationDetailUpdate{
		LastMessStatus: updateData.LastMessStatus,
	}

	return conversationDetailUpdate, nil
}
