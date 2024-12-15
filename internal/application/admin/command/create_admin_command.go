package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"time"
)

type CreateAdminCommand struct {
	ID          uuid.UUID
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
