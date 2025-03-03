package implement

import (
	"context"

	"github.com/google/uuid"
	adminCommand "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	adminEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	adminRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type sAdminInfo struct {
	adminRepo adminRepo.IAdminRepository
}

func NewAdminInfoImplement(
	adminRepo adminRepo.IAdminRepository,
) *sAdminInfo {
	return &sAdminInfo{
		adminRepo: adminRepo,
	}
}

func (s *sAdminInfo) UpdateAdmin(
	ctx context.Context,
	command *adminCommand.UpdateAdminInfoCommand,
) (result *adminCommand.UpdateAdminInfoCommandResult, err error) {
	// 1. Update admin info
	updateAdminEntity := &adminEntity.AdminUpdate{
		FamilyName:  command.FamilyName,
		Name:        command.Name,
		PhoneNumber: command.PhoneNumber,
		IdentityId:  command.IdentityId,
		Birthday:    command.Birthday,
	}

	if err = updateAdminEntity.ValidateAdminUpdate(); err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	adminFound, err := s.adminRepo.UpdateOne(ctx, *command.AdminID, updateAdminEntity)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if adminFound == nil {
		return nil, response.NewDataNotFoundError("admin not found")
	}

	// 2. Map to result
	return &adminCommand.UpdateAdminInfoCommandResult{
		Admin: mapper.NewAdminResult(adminFound),
	}, nil
}

func (s *sAdminInfo) GetAdminStatusById(
	ctx context.Context,
	id uuid.UUID,
) (status bool, err error) {
	adminStatus, err := s.adminRepo.GetStatusById(ctx, id)
	if err != nil {
		return false, err
	}

	return adminStatus, nil
}
