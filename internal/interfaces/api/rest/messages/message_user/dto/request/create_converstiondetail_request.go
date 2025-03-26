package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
)

type CreateConversationDetailRequest struct {
	UserId         uuid.UUID  `json:"user_id"`
	ConversationId uuid.UUID  `json:"conversation_id"`
	LastMessId     *uuid.UUID `json:"last_mess_id"`
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
	// Nếu LastMessId là nil, gán giá trị mặc định (UUID rỗng)
	lastMessId := uuid.Nil
	if req.LastMessId != nil {
		lastMessId = *req.LastMessId
	}

	return &command.CreateConversationDetailCommand{
		UserId:         req.UserId,
		ConversationId: req.ConversationId,
		LastMessId:     lastMessId,
	}, nil
}
