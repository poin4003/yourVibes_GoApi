package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	admin_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sSuperAdmin struct {
	adminRepo admin_repo.IAdminRepository
}

func NewSuperAdminImplement(
	adminRepo admin_repo.IAdminRepository,
) *sSuperAdmin {
	return &sSuperAdmin{
		adminRepo: adminRepo,
	}
}

func (s *sSuperAdmin) CreateAdmin(
	ctx context.Context,
	command *command.CreateAdminCommand,
) (result *command.CreateAdminCommandResult, err error) {
	return nil, nil
}

func (s *sSuperAdmin) UpdateAdmin(
	ctx context.Context,
	command *command.UpdateAdminForSuperAdminCommand,
) (result *command.UpdateAdminForSuperAdminCommandResult, err error) {
	return nil, nil
}
