package query

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
)

type ConversationDetailObject struct {
	// UserId         int `form:"user_id,omitempty"`
	ConversationId uuid.UUID `form:"conversation_id,omitempty"`
	Limit          int       `form:"limit,omitempty"`
	Page           int       `form:"page,omitempty"`
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

func (req *ConversationDetailObject) ToGetConversationDetailQuery(
	userId uuid.UUID,
	conversationId uuid.UUID,
) (*query.GetConversationDetailQuery, error) {
	return &query.GetConversationDetailQuery{
		UserId:         userId,
		ConversationId: conversationId,
		Limit:          req.Limit,
		Page:           req.Page,
	}, nil
}
