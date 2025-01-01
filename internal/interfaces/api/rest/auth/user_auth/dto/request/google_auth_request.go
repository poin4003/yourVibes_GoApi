package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type AuthGoogleRequest struct {
	OpenId       string          `json:"open_id"`
	AuthGoogleId string          `json:"auth_google_id"`
	Platform     consts.Platform `json:"platform"`
	FamilyName   string          `json:"family_name"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	AvatarUrl    string          `json:"avatar_url"`
}

func ValidateAuthGoogleRequest(req interface{}) error {
	dto, ok := req.(*AuthGoogleRequest)
	if !ok {
		return fmt.Errorf("input is not AuthGoogleRequest")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.OpenId, validation.Required),
		validation.Field(&dto.AuthGoogleId, validation.Required),
		validation.Field(&dto.Platform, validation.Required, validation.In(consts.WEB, consts.ANDROID, consts.IOS)),
		validation.Field(&dto.FamilyName, validation.Required, validation.Length(2, 255)),
		validation.Field(&dto.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&dto.Email, validation.Required, is.Email),
		validation.Field(&dto.AvatarUrl, validation.Required, validation.Length(10, 255)),
	)
}

func (req *AuthGoogleRequest) ToAuthGoogleCommand() (*user_command.AuthGoogleCommand, error) {
	return &user_command.AuthGoogleCommand{
		OpenId:       req.OpenId,
		AuthGoogleId: req.AuthGoogleId,
		Platform:     req.Platform,
		FamilyName:   req.FamilyName,
		Name:         req.Name,
		Email:        req.Email,
		AvatarUrl:    req.AvatarUrl,
	}, nil
}
