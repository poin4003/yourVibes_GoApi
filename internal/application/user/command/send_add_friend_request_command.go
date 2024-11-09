package command

import "github.com/google/uuid"

type SendAddFriendRequestCommand struct {
	UserId   uuid.UUID
	FriendId uuid.UUID
}

type SendAddFriendRequestCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
