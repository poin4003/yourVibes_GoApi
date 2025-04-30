package command

import "github.com/google/uuid"

type DeleteMessageCommand struct {
	MessageId           *uuid.UUID
	AuthenticatedUserId uuid.UUID
}
