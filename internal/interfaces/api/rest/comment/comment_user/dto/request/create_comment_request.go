package request

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
)

type CreateCommentInput struct {
	PostId   uuid.UUID  `json:"post_id" binding:"required"`
	ParentId *uuid.UUID `json:"parent_id,omitempty"`
	Content  string     `json:"content" binding:"required"`
}

func (req *CreateCommentInput) ToCreateCommentCommand(
	userId uuid.UUID,
) (*command.CreateCommentCommand, error) {
	return &command.CreateCommentCommand{
		PostId:   req.PostId,
		UserId:   userId,
		ParentId: req.ParentId,
		Content:  req.Content,
	}, nil
}
