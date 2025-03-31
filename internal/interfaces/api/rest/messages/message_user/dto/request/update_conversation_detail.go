package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
)

type UpdateConversationDetail struct {
	ConversationId uuid.UUID `json:"conversation_id"`
	UserId         uuid.UUID `json:"user_id"`
}

func ValidateUpdateConversationDetail(req interface{}) error {
	dto, ok := req.(*UpdateConversationDetail)

	if !ok {
		return fmt.Errorf("input is not UpdateConversationDetail")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.UserId, validation.Required),
		validation.Field(&dto.ConversationId, validation.Required),
	)
}

func (req *UpdateConversationDetail) ToUpdateConversationDetailCommand(
	userId uuid.UUID,
	conversationId uuid.UUID,
) (*command.UpdateOneStatusConversationDetailCommand, error) {
	return &command.UpdateOneStatusConversationDetailCommand{
		UserId:         userId,
		ConversationId: conversationId,
	}, nil
}
