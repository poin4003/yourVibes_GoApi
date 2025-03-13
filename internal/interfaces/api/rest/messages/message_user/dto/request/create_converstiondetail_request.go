package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
)

type CreateConversationDetailRequest struct {
	UserId         uuid.UUID `json:"user_id"`
	ConversationId uuid.UUID `json:"conversation_id"`
}

func ValidateCreatCOnversationDetailRequest(req interface{}) error {
	dto, ok := req.(*CreateConversationDetailRequest)

	if !ok {
		return fmt.Errorf("input is not CreatCOnversationDetailRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.UserId, validation.Required),
		validation.Field(&dto.ConversationId, validation.Required),
	)
}

func (req *CreateConversationDetailRequest) ToCreateConversationDetailCommand(
	userId uuid.UUID,
	conversationId uuid.UUID,
) (*command.CreateConversationDetailCommand, error) {
	return &command.CreateConversationDetailCommand{
		UserId:         req.UserId,
		ConversationId: req.ConversationId,
	}, nil
}
