package command

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/common"
)

type UpdateAdminInfoCommand struct {
	AdminID     *uuid.UUID
	FamilyName  *string
	Name        *string
	PhoneNumber *string
	IdentityId  *string
	Birthday    *time.Time
}

type UpdateAdminInfoCommandResult struct {
	Admin *common.AdminResult
}
