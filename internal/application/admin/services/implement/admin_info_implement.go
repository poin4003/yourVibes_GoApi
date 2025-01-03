package implement

import (
	"context"
	"errors"
	"github.com/google/uuid"
	adminCommand "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	adminEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	adminRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
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
	result = &adminCommand.UpdateAdminInfoCommandResult{}
	result.Admin = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Update admin info
	updateAdminEntity := &adminEntity.AdminUpdate{
		FamilyName:  command.FamilyName,
		Name:        command.Name,
		PhoneNumber: command.PhoneNumber,
		IdentityId:  command.IdentityId,
		Birthday:    command.Birthday,
	}

	if err = updateAdminEntity.ValidateAdminUpdate(); err != nil {
		return result, err
	}

	adminFound, err := s.adminRepo.UpdateOne(ctx, *command.AdminID, updateAdminEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			result.Admin = nil
			return result, err
		}
		return result, err
	}

	// 2. Map to result
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.Admin = mapper.NewAdminResult(adminFound)
	return result, nil
}

func (s *sAdminInfo) GetAdminStatusById(
	ctx context.Context,
	id uuid.UUID,
) (status bool, err error) {
	adminStatus, err := s.adminRepo.GetStatusById(ctx, id)
	if err != nil {
		return false, err
	}

	return adminStatus, err
}
