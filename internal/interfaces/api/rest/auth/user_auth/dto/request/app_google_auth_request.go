package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type AppAuthGoogleRequest struct {
	OpenId      string          `json:"open_id"`
	Platform    consts.Platform `json:"platform"`
	RedirectUrl string          `json:"redirect_url"`
}

func ValidateAppAuthGoogleRequest(req interface{}) error {
	dto, ok := req.(*AppAuthGoogleRequest)
	if !ok {
		return fmt.Errorf("input is not AuthGoogleRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.OpenId, validation.Required),
		validation.Field(&dto.Platform, validation.Required, validation.In(consts.Platforms...)),
		validation.Field(&dto.RedirectUrl, validation.Required),
	)
}

func (req *AppAuthGoogleRequest) ToAppAuthGoogleCommand() (*userCommand.AuthAppGoogleCommand, error) {
	return &userCommand.AuthAppGoogleCommand{
		OpenId:      req.OpenId,
		Platform:    req.Platform,
		RedirectUrl: req.RedirectUrl,
	}, nil
}
