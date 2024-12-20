package command

import "github.com/google/uuid"

type ActivateUserAccountCommand struct {
	UserId uuid.UUID
}

type ActivateUserAccountCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
