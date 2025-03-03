package command

import "github.com/google/uuid"

type DeletePostCommand struct {
	PostId *uuid.UUID
}
