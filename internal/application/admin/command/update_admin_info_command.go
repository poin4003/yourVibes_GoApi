package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"time"
)

type UpdateAdminInfoCommand struct {
	AdminID     *uuid.UUID
	FamilyName  *string
	Name        *string
	Email       *string
	Password    *string
	PhoneNumber *string
	IdentityId  *string
	Birthday    *time.Time
}

type UpdateAdminInfoCommandResult struct {
	Admin          *common.AdminResult
	ResultCode     int
	HttpStatusCode int
}