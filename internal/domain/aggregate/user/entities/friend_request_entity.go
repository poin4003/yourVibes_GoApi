package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type FriendRequest struct {
	UserId   uuid.UUID `validate:"required,uuid4"`
	FriendId uuid.UUID `validate:"required,uuid4"`
}

func (fr *FriendRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(fr)
}

func NewFriendRequest(
	userId uuid.UUID,
	friendId uuid.UUID,
) (*FriendRequest, error) {
	newFriendRequest := &FriendRequest{
		UserId:   userId,
		FriendId: friendId,
	}
	if err := newFriendRequest.Validate(); err != nil {
		return nil, err
	}

	return newFriendRequest, nil
}
