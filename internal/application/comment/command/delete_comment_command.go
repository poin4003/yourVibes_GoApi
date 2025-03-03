package command

import "github.com/google/uuid"

type DeleteCommentCommand struct {
	CommentId uuid.UUID
}
