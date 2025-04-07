package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/google/uuid"
	adminCommand "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	adminEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	repository "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sAdminInfo struct {
	adminRepo  repository.IAdminRepository
	adminCache cache.IAdminCache
}

func NewAdminInfoImplement(
	adminRepo repository.IAdminRepository,
	adminCache cache.IAdminCache,
) *sAdminInfo {
	return &sAdminInfo{
		adminRepo:  adminRepo,
		adminCache: adminCache,
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
) (status *bool, err error) {
	// 1. Get admin status from cache
	adminStatus := s.adminCache.GetAdminStatus(ctx, id)
	// 2. Check if cache miss
	if adminStatus == nil {
		adminStatus, err = s.adminRepo.GetStatusById(ctx, id)
		if err != nil {
			return nil, err
		}
		go func(adminId uuid.UUID, adminStatus bool) {
			s.adminCache.SetAdminStatus(ctx, adminId, adminStatus)
		}(id, *adminStatus)
	}

	return adminStatus, nil
}
