package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	user_validator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/validator"
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
	user user_validator.ValidatedUser,
	friend user_validator.ValidatedUser,
) (*Friend, error) {
	newFriend := &Friend{
		UserId:   userId,
		FriendId: friendId,
		User:     user.User,
		Friend:   friend.User,
	}
	if err := friend.Validate(); err != nil {
		return nil, err
	}

	return newFriend, nil
}
