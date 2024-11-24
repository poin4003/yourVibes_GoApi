package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type Friend struct {
	UserId   uuid.UUID `validate:"required,uuid4"`
	FriendId uuid.UUID `validate:"required,uuid4"`
}

func (f *Friend) ValidateFriend() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.UserId, validation.Required),
		validation.Field(&f.FriendId, validation.Required),
	)
}

func NewFriend(
	userId uuid.UUID,
	friendId uuid.UUID,
) (*Friend, error) {
	newFriend := &Friend{
		UserId:   userId,
		FriendId: friendId,
	}
	if err := newFriend.ValidateFriend(); err != nil {
		return nil, err
	}

	return newFriend, nil
}
