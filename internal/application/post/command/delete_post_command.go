package command

import "github.com/google/uuid"

type DeletePostCommand struct {
	PostId *uuid.UUID
}

type DeletePostCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
