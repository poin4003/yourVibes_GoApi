package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"

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
}

type RegisterCommandResult struct {
	User       *common.UserWithSettingResult
	ResultCode int
}

type LoginCommandResult struct {
	User           *common.UserWithSettingResult
	AccessToken    *string
	ResultCode     int
	HttpStatusCode int
}

type ChangePasswordCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}

type GetOtpForgotUserPasswordCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}

type ForgotUserPasswordCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}

type AuthGoogleCommandResult struct {
	User           *common.UserWithSettingResult
	AccessToken    *string
	ResultCode     int
	HttpStatusCode int
}
