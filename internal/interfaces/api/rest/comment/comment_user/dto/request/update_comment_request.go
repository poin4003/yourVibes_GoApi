package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"unicode/utf8"
)

type UpdateCommentRequest struct {
	Content *string `json:"content,omitempty"`
}

func ValidateUpdateCommentRequest(req interface{}) error {
	dto, ok := req.(*UpdateCommentRequest)
	if !ok {
		return fmt.Errorf("input is not UpdateCommentRequest")
	}

	return validation.ValidateStruct(dto,
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

func (req *UpdateCommentRequest) ToUpdateCommentCommand(
	commentId uuid.UUID,
) (*command.UpdateCommentCommand, error) {
	return &command.UpdateCommentCommand{
		CommentId: commentId,
		Content:   req.Content,
	}, nil
}
