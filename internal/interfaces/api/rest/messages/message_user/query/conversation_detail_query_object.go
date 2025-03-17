package query

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
)

type ConversationDetailObject struct {
	UserId         string `form:"user_id,omitempty"`
	ConversationId string `form:"conversation_id,omitempty"`
	Limit          int    `form:"limit,omitempty"`
	Page           int    `form:"page,omitempty"`
}

func ValidateConversationDetailObject(input interface{}) error {
	query, ok := input.(*ConversationDetailObject)
	if !ok {
		return fmt.Errorf("validateConversationDetailObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *ConversationDetailObject) ToGetConversationDetailQuery() (*query.GetConversationDetailQuery, error) {
	var userId uuid.UUID
	var conversationId uuid.UUID
	if req.UserId != "" {
		parseUserId, err := uuid.Parse(req.UserId)
		if err != nil {
			return nil, err
		}
		userId = parseUserId
	}

	if req.ConversationId != "" {
		parseConversationId, err := uuid.Parse(req.ConversationId)
		if err != nil {
			return nil, err
		}

		conversationId = parseConversationId
	}
	return &query.GetConversationDetailQuery{

		UserId:         userId,
		ConversationId: conversationId,
		Limit:          req.Limit,
		Page:           req.Page,
	}, nil
}
