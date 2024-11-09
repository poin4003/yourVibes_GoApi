package command

import "github.com/google/uuid"

type RemoveFriendRequestCommand struct {
	UserId   uuid.UUID
	FriendId uuid.UUID
}

type RemoveFriendRequestCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
