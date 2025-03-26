package request

import (
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
)

type CreateConversationRequest struct {
	Name    string               `form:"name"`
	Image   multipart.FileHeader `form:"image,omitempty" binding:"omitempty"`
	UserIds []string             `form:"user_ids" binding:"required"`
}

func ValidateCreateConversationRequest(req interface{}) error {
	dto, ok := req.(*CreateConversationRequest)

	if !ok {
		return fmt.Errorf("input is not CreateConversationRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Name, validation.Required, validation.Length(1, 30)),
		validation.Field(&dto.UserIds, validation.Required),
	)
}

func (req *CreateConversationRequest) ToCreateConversationCommand(
	Name string,
	userIds []uuid.UUID,
) *command.CreateConversationCommand {
	return &command.CreateConversationCommand{
		Name:    Name,
		UserIds: userIds,
	}
}
