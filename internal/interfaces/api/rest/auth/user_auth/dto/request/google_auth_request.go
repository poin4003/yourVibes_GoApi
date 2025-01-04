package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type AuthGoogleRequest struct {
	AuthorizationCode string          `json:"authorization_code"`
	Platform          consts.Platform `json:"platform"`
	RedirectUrl       string          `json:"redirect_url"`
}

func ValidateAuthGoogleRequest(req interface{}) error {
	dto, ok := req.(*AuthGoogleRequest)
	if !ok {
		return fmt.Errorf("input is not AuthGoogleRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.AuthorizationCode, validation.Required),
		validation.Field(&dto.Platform, validation.Required, validation.In(consts.WEB, consts.ANDROID, consts.IOS)),
		validation.Field(&dto.RedirectUrl, validation.Required),
	)
}

func (req *AuthGoogleRequest) ToAuthGoogleCommand() (*userCommand.AuthGoogleCommand, error) {
	return &userCommand.AuthGoogleCommand{
		AuthorizationCode: req.AuthorizationCode,
		Platform:          req.Platform,
		RedirectUrl:       req.RedirectUrl,
	}, nil
}
