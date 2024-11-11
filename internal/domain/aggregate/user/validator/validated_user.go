package validator

import (
	"fmt"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
)

type ValidatedUser struct {
	user_entity.User
	isValidated bool
}

func (vu *ValidatedUser) IsValid() bool {
	return vu.isValidated
}

func NewValidatedUser(user *user_entity.User) (*ValidatedUser, error) {
	if user == nil {
		return nil, fmt.Errorf("NewValidatedUser: user is nil")
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return &ValidatedUser{
		User:        *user,
		isValidated: true,
	}, nil
}
