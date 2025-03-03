package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
)

type CreateCommentCommand struct {
	PostId   uuid.UUID
	ParentId *uuid.UUID
	UserId   uuid.UUID
	Content  string
}

type CreateCommentResult struct {
	Comment *common.CommentResult
}
