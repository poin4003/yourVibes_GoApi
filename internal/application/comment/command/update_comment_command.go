package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
)

type UpdateCommentCommand struct {
	CommentId uuid.UUID
	Content   *string
}

type UpdateCommentResult struct {
	Comment        *common.CommentResult
	ResultCode     int
	HttpStatusCode int
}
