package query

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
)

type MessageObject struct {
	ConversationId string `form:"conversation_id,omitempty"`
	SortBy         string `form:"sort_by,omitempty"`
	IsDescending   bool   `form:"is_descending,omitempty"`
	Limit          int    `form:"limit,omitempty"`
	Page           int    `form:"page,omitempty"`
}

func ValidateMessageQueryObject(input interface{}) error {
	query, ok := input.(*MessageObject)
	if !ok {
		return fmt.Errorf("validateMessageQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.ConversationId, validation.Required),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *MessageObject) ToGetManyMessageQuery() (*query.GetMessagesByConversationIdQuery, error) {
	var conversationId uuid.UUID
	if req.ConversationId != "" {
		parseConversationId, err := uuid.Parse(req.ConversationId)
		if err != nil {
			return nil, err
		}
		conversationId = parseConversationId
	}
	return &query.GetMessagesByConversationIdQuery{
		ConversationId: conversationId,
		SortBy:         req.SortBy,
		IsDescending:   req.IsDescending,
		Limit:          req.Limit,
		Page:           req.Page,
	}, nil
}
