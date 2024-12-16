package implement

import (
	"context"
	"errors"
	"fmt"
	admin_command "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	admin_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	admin_validator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/validator"
	admin_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/crypto"
	"gorm.io/gorm"
	"net/http"
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
	command *admin_command.CreateAdminCommand,
) (result *admin_command.CreateAdminCommandResult, err error) {
	result = &admin_command.CreateAdminCommandResult{}
	// 1. Check admin exist
	adminFound, err := s.adminRepo.CheckAdminExistByEmail(ctx, command.Email)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.Admin = nil
	}

	if adminFound {
		result.ResultCode = response.ErrCodeAdminHasExist
		result.HttpStatusCode = http.StatusBadRequest
		result.Admin = nil
		return result, fmt.Errorf("admin %s already exists", command.Email)
	}

	// 2. Hash password
	hashedPassword, err := crypto.HashPassword(command.Password)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.Admin = nil
		return result, err
	}

	// 3. Create new admin
	newAdmin, err := admin_entity.NewAdmin(
		command.FamilyName,
		command.Name,
		command.Email,
		hashedPassword,
		command.PhoneNumber,
		command.IdentityId,
		command.Birthday,
		command.Role,
	)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	createdAdmin, err := s.adminRepo.CreateOne(ctx, newAdmin)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 4. Map to result
	validateAdmin, err := admin_validator.NewValidateAdmin(createdAdmin)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	result.Admin = mapper.NewAdminResultFromValidateEntity(validateAdmin)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sSuperAdmin) UpdateAdmin(
	ctx context.Context,
	command *admin_command.UpdateAdminForSuperAdminCommand,
) (result *admin_command.UpdateAdminForSuperAdminCommandResult, err error) {
	result = &admin_command.UpdateAdminForSuperAdminCommandResult{}
	// 1. Update admin status or role
	updateAdminEntity := &admin_entity.AdminUpdate{
		Status: command.Status,
		Role:   command.Role,
	}

	if err = updateAdminEntity.ValidateAdminUpdate(); err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.Admin = nil
		return result, err
	}

	adminFound, err := s.adminRepo.UpdateOne(ctx, *command.AdminId, updateAdminEntity)
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
