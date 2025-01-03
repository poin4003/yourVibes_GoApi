package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
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
		validation.Field(&dto.Content, validation.Required, validation.Length(2, 500)),
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
