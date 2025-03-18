package request

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"time"
)

type CreateMessageRequest struct {
	ConversationId uuid.UUID  `json:"conversation_id"`
	Content        string     `json:"content"`
	ParentId       *uuid.UUID `json:"parent_id,omitempty"`
	ParentContent  *string    `json:"parent_content,omitempty"`
	User           userDto    `json:"user"`
}

type userDto struct {
	ID         string `json:"id"`
	FamilyName string `json:"family_name"`
	Name       string `json:"name"`
	AvatarUrl  string `json:"avatar_url"`
}

func ValidateCreateMessageRequest(req interface{}) error {
	dto, ok := req.(*CreateMessageRequest)

	if !ok {
		return fmt.Errorf("input is not CreateMessageRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.ConversationId, validation.Required),
		validation.Field(&dto.Content, validation.Required, validation.Length(1, 500)),
		validation.Field(&dto.ParentContent, validation.By(func(value interface{}) error {
			if dto.ParentId != nil && value == nil {
				return errors.New("parent content is required when parent id provides")
			}
			return nil
		})),
		validation.Field(&dto.User, validation.Required),
	)
}

func (req *CreateMessageRequest) ToCreateMessageCommand(
	userId uuid.UUID,
) (*command.CreateMessageCommand, error) {
	user := command.UserCommand{
		ID:         req.User.ID,
		FamilyName: req.User.FamilyName,
		Name:       req.User.Name,
		AvatarUrl:  req.User.AvatarUrl,
	}

	return &command.CreateMessageCommand{
		ConversationId: req.ConversationId,
		UserId:         userId,
		Content:        req.Content,
		ParentId:       req.ParentId,
		ParentContent:  req.ParentContent,
		CreatedAt:      time.Now(),
		User:           user,
	}, nil
}
