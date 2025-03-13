package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
)

type CreateMessageRequest struct {
	ConversationId uuid.UUID  `json:"conversation_id"`
	Content        string     `json:"content"`
	ParentId       *uuid.UUID `json:"parent_id,omitempty"`
}

func ValidateCreateMessageRequest(req interface{}) error {
	dto, ok := req.(*CreateMessageRequest)

	if !ok {
		return fmt.Errorf("input is not CreateMessageRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.ConversationId, validation.Required),
		validation.Field(&dto.Content, validation.Required, validation.Length(1, 500)),
	)
}

func (req *CreateMessageRequest) ToCreateMessageCommand(
	userId uuid.UUID,
) (*command.CreateMessageCommand, error) {
	return &command.CreateMessageCommand{
		ConversationId: req.ConversationId,
		UserId:         userId,
		Content:        req.Content,
		ParentId:       req.ParentId,
	}, nil
}
