package request

import (
	"fmt"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
)

type CreateConversationRequest struct {
	Name  string               `json:"name"`
	Image multipart.FileHeader `form:"image,omitempty" binding:"omitempty"`
}

func ValidateCreateConversationRequest(req interface{}) error {
	dto, ok := req.(*CreateConversationRequest)

	if !ok {
		return fmt.Errorf("input is not CreateConversationRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Name, validation.Required, validation.Length(1, 30)),
	)
}

func (req *CreateConversationRequest) ToCreateConversationCommand(
	Name string) *command.CreateConversationCommand {
	return &command.CreateConversationCommand{
		Name: Name,
	}
}
