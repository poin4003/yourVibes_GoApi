package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
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
		validation.Field(&dto.Content, validation.RuneLength(1, 500)),
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
