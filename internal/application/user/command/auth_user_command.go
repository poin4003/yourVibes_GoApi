package command

import (
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
