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

type AuthGoogleCommand struct {
	OpenId       string
	AuthGoogleId string
	Platform     consts.Platform
	FamilyName   string
	Name         string
	Email        string
	AvatarUrl    string
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

type AuthGoogleCommandResult struct {
	User           *common.UserWithSettingResult
	AccessToken    *string
	ResultCode     int
	HttpStatusCode int
}
