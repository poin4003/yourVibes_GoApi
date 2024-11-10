package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Friend struct {
	UserId   uuid.UUID `validate:"required,uuid4"`
	FriendId uuid.UUID `validate:"required,uuid4"`
	User     User      `validate:"required"`
	Friend   User      `validate:"required"`
}

func (f *Friend) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}

func NewFriend(
	userId uuid.UUID,
	friendId uuid.UUID,
) (*Friend, error) {
	newFriend := &Friend{
		UserId:   userId,
		FriendId: friendId,
	}
	if err := newFriend.Validate(); err != nil {
		return nil, err
	}

	return newFriend, nil
}
