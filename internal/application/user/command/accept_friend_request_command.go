package command

import "github.com/google/uuid"

type AcceptFriendRequestCommand struct {
	UserId   uuid.UUID
	FriendId uuid.UUID
}

type AcceptFriendRequestCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
