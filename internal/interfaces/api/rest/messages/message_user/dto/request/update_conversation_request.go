package request

import (
	"fmt"
	"mime/multipart"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
)

type UpdateConversationRequest struct {
	Name  *string              `form:"name,omitempty"`
	Image multipart.FileHeader `form:"image,omitempty" binding:"omitempty"`
}

func ValidateUpdateConversationRequest(req interface{}) error {
	dto, ok := req.(*UpdateConversationRequest)
	if !ok {
		return fmt.Errorf("validate UpdateConversationRequest failed")
	}
	return validation.ValidateStruct(dto,
		validation.Field(&dto.Name, validation.RuneLength(1, 30)),
		validation.Field(&dto.Image, validation.By(validateImage)),
	)
}

func validateImage(value interface{}) error {
	if value == nil {
		return nil
	}
	fileHeader, ok := value.(multipart.FileHeader)
	if !ok {
		return nil
	}

	if fileHeader.Size == 0 {
		return nil
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("file must be an image")
	}
	return nil
}

func (req *UpdateConversationRequest) ToUpdateConversationCommand(
	conversationId uuid.UUID,
	image *multipart.FileHeader,
) (*command.UpdateConversationCommand, error) {
	return &command.UpdateConversationCommand{
		ConversationId: &conversationId,
		Name:           req.Name,
		Image:          image,
	}, nil
}
