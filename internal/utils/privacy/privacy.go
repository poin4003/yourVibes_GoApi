package privacy

import (
	"github.com/go-playground/validator/v10"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

func IsValid(p consts.PrivacyLevel) bool {
	switch p {
	case consts.PUBLIC, consts.PRIVATE, consts.FRIEND_ONLY:
		return true
	}
	return false
}

func ValidatePrivacy(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == string(consts.PUBLIC) || value == string(consts.PRIVATE) || value == string(consts.FRIEND_ONLY)
}
