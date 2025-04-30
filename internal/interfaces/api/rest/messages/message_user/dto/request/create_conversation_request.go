package request

import (
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"

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
		validation.Field(&dto.Name, validation.Required, validation.RuneLength(1, 30)),
		validation.Field(&dto.UserIds,
			validation.Required,
			validation.By(checkDuplicateUserIds),
		),
	)
}

func checkDuplicateUserIds(value interface{}) error {
	userIds, ok := value.([]string)
	if !ok {
		return fmt.Errorf("invalid type for user_ids")
	}

	seen := make(map[string]bool)
	for _, id := range userIds {
		if seen[id] {
			return fmt.Errorf("duplicate user id found: %s", id)
		}
		seen[id] = true
	}

	return nil
}

func (req *CreateConversationRequest) ToCreateConversationCommand(
	Name string,
	userIds []uuid.UUID,
	ownerId uuid.UUID,
) *command.CreateConversationCommand {
	return &command.CreateConversationCommand{
		Name:    Name,
		UserIds: userIds,
		OwnerId: ownerId,
	}
}
