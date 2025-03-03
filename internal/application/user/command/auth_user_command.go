package command

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"

	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
)

type RegisterCommand struct {
	FamilyName  string
	Name        string
	Email       string
	Password    string
	PhoneNumber string
	Birthday    time.Time
	Otp         string
}

type LoginCommand struct {
	Email    string
	Password string
}

type ChangePasswordCommand struct {
	UserId      uuid.UUID
	OldPassword string
	NewPassword string
}

type GetOtpForgotUserPasswordCommand struct {
	Email string
}

type ForgotUserPasswordCommand struct {
	Email       string
	Otp         string
	NewPassword string
}

type AuthGoogleCommand struct {
	AuthorizationCode string
	Platform          consts.Platform
	RedirectUrl       string
}

type AuthAppGoogleCommand struct {
	OpenId      string
	Platform    consts.Platform
	RedirectUrl string
}

type RegisterCommandResult struct {
	User *common.UserWithSettingResult
}

type LoginCommandResult struct {
	User        *common.UserWithSettingResult
	AccessToken *string
}

type AuthGoogleCommandResult struct {
	User        *common.UserWithSettingResult
	AccessToken *string
}

type AuthAppGoogleCommandResult struct {
	User        *common.UserWithSettingResult
	AccessToken *string
}
