package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type FriendRequest struct {
	UserId   uuid.UUID
	FriendId uuid.UUID
}

func (fr *FriendRequest) ValidateFriend() error {
	return validation.ValidateStruct(fr,
		validation.Field(&fr.FriendId, validation.Required),
		validation.Field(&fr.UserId, validation.Required),
	)
}

func NewFriendRequest(
	userId uuid.UUID,
	friendId uuid.UUID,
) (*FriendRequest, error) {
	newFriendRequest := &FriendRequest{
		UserId:   userId,
		FriendId: friendId,
	}
	if err := newFriendRequest.ValidateFriend(); err != nil {
		return nil, err
	}

	return newFriendRequest, nil
}
