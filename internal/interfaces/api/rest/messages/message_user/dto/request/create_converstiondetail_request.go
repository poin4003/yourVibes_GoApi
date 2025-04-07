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

type CreateManyConversationDetailRequest struct {
	UserIds        []string  `json:"user_ids"`
	ConversationId uuid.UUID `json:"conversation_id"`
}

func ValidateCreateConversationDetailRequest(req interface{}) error {
	dto, ok := req.(*CreateConversationDetailRequest)

	if !ok {
		return fmt.Errorf("input is not CreatCOnversationDetailRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.UserId, validation.Required),
		validation.Field(&dto.ConversationId, validation.Required),
	)
}

func ValidateCreateManyConversationDetailRequest(req interface{}) error {
	dto, ok := req.(*CreateManyConversationDetailRequest)

	if !ok {
		return fmt.Errorf("input is not CreatCOnversationDetailRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.UserIds, validation.Required, validation.Length(1, 0)), // Phải có ít nhất 1 user
		validation.Field(&dto.ConversationId, validation.Required),
	)
}

func (req *CreateConversationDetailRequest) ToCreateConversationDetailCommand() (*command.CreateConversationDetailCommand, error) {

	return &command.CreateConversationDetailCommand{
		UserId:         req.UserId,
		ConversationId: req.ConversationId,
	}, nil
}
func (req *CreateManyConversationDetailRequest) ToCreateManyConversationDetailCommands(
	userIds []uuid.UUID,
) (*command.CreateManyConversationDetailCommand, error) {
	if len(req.UserIds) == 0 {
		return nil, fmt.Errorf("user_ids is empty")
	}
	return &command.CreateManyConversationDetailCommand{
		UserIds:        userIds,
		ConversationId: req.ConversationId,
	}, nil

}
