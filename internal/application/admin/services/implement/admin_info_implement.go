package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	admin_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sAdminInfo struct {
	adminRepo admin_repo.IAdminRepository
}

func NewAdminInfoImplement(
	adminRepo admin_repo.IAdminRepository,
) *sAdminInfo {
	return &sAdminInfo{
		adminRepo: adminRepo,
	}
}

func (s *sAdminInfo) UpdateAdmin(
	ctx context.Context,
	command *command.UpdateAdminInfoCommand,
) (result *command.UpdateAdminInfoCommandResult, err error) {
	return nil, nil
}
