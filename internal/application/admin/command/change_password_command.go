package command

import "github.com/google/uuid"

type ChangeAdminPasswordCommand struct {
	AdminId     uuid.UUID
	OldPassword string
	NewPassword string
}

type ChangeAdminPasswordCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
