package command

import "github.com/google/uuid"

type ActivateUserAccountCommand struct {
	UserId uuid.UUID
}
