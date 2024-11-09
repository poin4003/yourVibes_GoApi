package command

import (
	"github.com/google/uuid"
)

type UnFriendCommand struct {
	UserId   uuid.UUID
	FriendId uuid.UUID
}

type UnFriendCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
