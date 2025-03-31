package query

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	conversationQuery "github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
)

type ConversationObject struct {
	Name         string    `form:"name,omitempty"`
	CreatedAt    time.Time `form:"created_at,omitempty"`
	SortBy       string    `form:"sort_by,omitempty"`
	IsDescending bool      `form:"isDescending,omitempty"`
	Limit        int       `form:"limit,omitempty"`
	Page         int       `form:"page,omitempty"`
}

func ValidateConversationObject(input interface{}) error {
	query, ok := input.(*ConversationObject)
	if !ok {
		return fmt.Errorf("validateConversationObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *ConversationObject) ToGetManyConversationQuery() (*conversationQuery.GetManyConversationQuery, error) {

	return &conversationQuery.GetManyConversationQuery{
		Name:         req.Name,
		CreatedAt:    req.CreatedAt,
		SortBy:       req.SortBy,
		IsDescending: req.IsDescending,
		Limit:        req.Limit,
		Page:         req.Page,
	}, nil
}

func (req *ConversationObject) ToGetOneConversationQuery(
	conversationId uuid.UUID,
	authenticatedUserId uuid.UUID,
) (*conversationQuery.GetOneConversationQuery, error) {
	return &conversationQuery.GetOneConversationQuery{
		ConversationId: conversationId,
	}, nil
}
