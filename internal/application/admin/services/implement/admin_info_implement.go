package implement

import (
	"context"
	"errors"
	"github.com/google/uuid"
	admin_command "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	admin_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	admin_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
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
	command *admin_command.UpdateAdminInfoCommand,
) (result *admin_command.UpdateAdminInfoCommandResult, err error) {
	result = &admin_command.UpdateAdminInfoCommandResult{}
	// 1. Update admin info
	updateAdminEntity := &admin_entity.AdminUpdate{
		FamilyName:  command.FamilyName,
		Name:        command.Name,
		PhoneNumber: command.PhoneNumber,
		IdentityId:  command.IdentityId,
		Birthday:    command.Birthday,
	}

	if err = updateAdminEntity.ValidateAdminUpdate(); err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.Admin = nil
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
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.Admin = nil
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
