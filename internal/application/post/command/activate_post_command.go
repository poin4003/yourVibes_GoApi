package command

import "github.com/google/uuid"

type ActivatePostCommand struct {
	PostId uuid.UUID
}

type ActivatePostCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
