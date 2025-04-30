package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type ConversationDetail struct {
	UserId           uuid.UUID
	ConversationId   uuid.UUID
	User             *User
	Conversation     *Conversation
	LastMessStatus   bool
	LastMess         *string
	ConversationRole consts.ConversationRole
}

type ConversationDetailUpdate struct {
	LastMessStatus *bool
}

func (cd *ConversationDetail) Validate() error {
	return validation.ValidateStruct(cd,
		validation.Field(&cd.UserId, validation.Required),
		validation.Field(&cd.ConversationId, validation.Required),
		validation.Field(&cd.LastMessStatus, validation.Required),
		validation.Field(&cd.ConversationRole, validation.In(consts.ConversationRoles...)),
	)
}

func NewConversationDetail(
	userId uuid.UUID,
	conversationId uuid.UUID,
	conversationRole consts.ConversationRole,
) (*ConversationDetail, error) {
	conversationDetail := &ConversationDetail{
		UserId:           userId,
		ConversationId:   conversationId,
		LastMessStatus:   true,
		LastMess:         nil,
		ConversationRole: conversationRole,
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
