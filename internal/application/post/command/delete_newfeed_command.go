package command

import (
	"github.com/google/uuid"
)

type DeleteNewFeedCommand struct {
	UserId uuid.UUID
	PostId uuid.UUID
}
