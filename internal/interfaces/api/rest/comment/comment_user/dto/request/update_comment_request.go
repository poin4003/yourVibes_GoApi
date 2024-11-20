package request

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
)

type UpdateCommentInput struct {
	Content *string `json:"content,omitempty"`
}

func (req *UpdateCommentInput) ToUpdateCommentCommand(
	commentId uuid.UUID,
) (*command.UpdateCommentCommand, error) {
	return &command.UpdateCommentCommand{
		CommentId: commentId,
		Content:   req.Content,
	}, nil
}
