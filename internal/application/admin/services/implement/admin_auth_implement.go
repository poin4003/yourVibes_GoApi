package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	admin_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sAdminAuth struct {
	adminRepo admin_repo.IAdminRepository
}

func NewAdminAuthImplement(
	adminRepo admin_repo.IAdminRepository,
) *sAdminAuth {
	return &sAdminAuth{
		adminRepo: adminRepo,
	}
}

func (s *sAdminAuth) Login(
	ctx context.Context,
	command *command.LoginCommand,
) (result *command.LoginCommandResult, err error) {
	return nil, nil
}
