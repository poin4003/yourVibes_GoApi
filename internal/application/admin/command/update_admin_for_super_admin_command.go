package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
)

type UpdateAdminForSuperAdminCommand struct {
	AdminId *uuid.UUID
	Role    *bool
	Status  *bool
}

type UpdateAdminForSuperAdminCommandResult struct {
	Admin          *common.AdminResult
	ResultCode     int
	HttpStatusCode int
}
