package implement

import (
	"context"
	"errors"
	"fmt"
	adminCommand "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/query"
	adminQuery "github.com/poin4003/yourVibes_GoApi/internal/application/admin/query"
	adminEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	adminValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/validator"
	adminRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/crypto"
	"gorm.io/gorm"
	"net/http"
)

type sSuperAdmin struct {
	adminRepo adminRepo.IAdminRepository
}

func NewSuperAdminImplement(
	adminRepo adminRepo.IAdminRepository,
) *sSuperAdmin {
	return &sSuperAdmin{
		adminRepo: adminRepo,
	}
}

func (s *sSuperAdmin) CreateAdmin(
	ctx context.Context,
	command *adminCommand.CreateAdminCommand,
) (result *adminCommand.CreateAdminCommandResult, err error) {
	result = &adminCommand.CreateAdminCommandResult{}
	result.Admin = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Check admin exist
	adminFound, err := s.adminRepo.CheckAdminExistByEmail(ctx, command.Email)
	if err != nil {
		return result, err
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
		return result, err
	}

	// 3. Create new admin
	newAdmin, err := adminEntity.NewAdmin(
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
		return result, err
	}

	createdAdmin, err := s.adminRepo.CreateOne(ctx, newAdmin)
	if err != nil {
		return result, err
	}

	// 4. Map to result
	validateAdmin, err := adminValidator.NewValidateAdmin(createdAdmin)
	if err != nil {
		return result, err
	}

	result.Admin = mapper.NewAdminResultFromValidateEntity(validateAdmin)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sSuperAdmin) UpdateAdmin(
	ctx context.Context,
	command *adminCommand.UpdateAdminForSuperAdminCommand,
) (result *adminCommand.UpdateAdminForSuperAdminCommandResult, err error) {
	result = &adminCommand.UpdateAdminForSuperAdminCommandResult{}
	result.Admin = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Update admin status or role
	updateAdminEntity := &adminEntity.AdminUpdate{
		Status: command.Status,
		Role:   command.Role,
	}

	if err = updateAdminEntity.ValidateAdminUpdate(); err != nil {
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
		return result, err
	}

	// 2. Map to result
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.Admin = mapper.NewAdminResult(adminFound)
	return result, nil
}

func (s *sSuperAdmin) GetOneAdmin(
	ctx context.Context,
	query *query.GetOneAdminQuery,
) (result *query.AdminQueryResult, err error) {
	result = &adminQuery.AdminQueryResult{}
	result.Admin = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Get admin info
	adminFound, err := s.adminRepo.GetById(ctx, query.AdminId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			result.Admin = nil
			return result, nil
		}
		return result, err
	}

	// 2. Map to result
	result.Admin = mapper.NewAdminResult(adminFound)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sSuperAdmin) GetManyAdmin(
	ctx context.Context,
	query *query.GetManyAdminQuery,
) (result *query.AdminQueryListResult, err error) {
	result = &adminQuery.AdminQueryListResult{}
	result.Admins = nil
	result.PagingResponse = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Get list admin
	adminEntities, paging, err := s.adminRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	var adminResults []*common.AdminResult
	for _, admin := range adminEntities {
		adminResult := mapper.NewAdminResult(admin)
		adminResults = append(adminResults, adminResult)
	}

	result.Admins = adminResults
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}
