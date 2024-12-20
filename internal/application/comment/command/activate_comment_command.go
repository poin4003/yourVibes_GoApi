package command

import "github.com/google/uuid"

type ActivateCommentCommand struct {
	CommentId uuid.UUID
}

type ActivateCommentCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
