package command

import "github.com/google/uuid"

type RemoveFriendRequestCommand struct {
	UserId   uuid.UUID
	FriendId uuid.UUID
}
