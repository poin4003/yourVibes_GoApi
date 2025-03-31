package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"unicode/utf8"
)

type CreateCommentRequest struct {
	PostId   uuid.UUID  `json:"post_id"`
	ParentId *uuid.UUID `json:"parent_id,omitempty"`
	Content  string     `json:"content"`
}

func ValidateCreateCommentRequest(req interface{}) error {
	dto, ok := req.(*CreateCommentRequest)
	if !ok {
		return fmt.Errorf("input is not CreateCommentRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.PostId, validation.Required),
		validation.Field(&dto.Content, validation.By(func(value interface{}) error {
			str, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid content type")
			}

			length := utf8.RuneCountInString(str)
			if length < 1 || length > 500 {
				return fmt.Errorf("content length must be between 2 and 500 characters, but got %d", length)
			}
			return nil
		})),
	)
}

func (req *CreateCommentRequest) ToCreateCommentCommand(
	userId uuid.UUID,
) (*command.CreateCommentCommand, error) {
	return &command.CreateCommentCommand{
		PostId:   req.PostId,
		UserId:   userId,
		ParentId: req.ParentId,
		Content:  req.Content,
	}, nil
}
