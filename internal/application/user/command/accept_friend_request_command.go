package command

import "github.com/google/uuid"

type AcceptFriendRequestCommand struct {
	UserId   uuid.UUID
	FriendId uuid.UUID
}
