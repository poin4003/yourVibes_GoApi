package command

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/common"
	"time"
)

type CreateAdminCommand struct {
	FamilyName  string
	Name        string
	Email       string
	Password    string
	PhoneNumber string
	IdentityId  string
	Birthday    time.Time
	Role        bool
}

type CreateAdminCommandResult struct {
	Admin          *common.AdminResult
	ResultCode     int
	HttpStatusCode int
}
