package command

import "github.com/google/uuid"

type DeleteCommentCommand struct {
	CommentId uuid.UUID
}

type DeleteCommentResult struct {
	ResultCode     int
	HttpStatusCode int
}
