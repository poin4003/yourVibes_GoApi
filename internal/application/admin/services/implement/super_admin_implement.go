package implement

import (
	"context"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/crypto"

	adminCommand "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	adminQuery "github.com/poin4003/yourVibes_GoApi/internal/application/admin/query"
	adminEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	adminValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/validator"
	adminRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
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
	// 1. Check admin exist
	adminFound, err := s.adminRepo.CheckAdminExistByEmail(ctx, command.Email)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	if adminFound {
		return nil, response2.NewCustomError(
			response2.ErrDataHasAlreadyExist,
			"admin already exist",
		)
	}

	// 2. Hash password
	hashedPassword, err := crypto.HashPassword(command.Password)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
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
		return nil, response2.NewServerFailedError(err.Error())
	}

	createdAdmin, err := s.adminRepo.CreateOne(ctx, newAdmin)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	// 4. Map to result
	validateAdmin, err := adminValidator.NewValidateAdmin(createdAdmin)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	return &adminCommand.CreateAdminCommandResult{
		Admin: mapper.NewAdminResultFromValidateEntity(validateAdmin),
	}, nil
}

func (s *sSuperAdmin) UpdateAdmin(
	ctx context.Context,
	command *adminCommand.UpdateAdminForSuperAdminCommand,
) (result *adminCommand.UpdateAdminForSuperAdminCommandResult, err error) {
	adminFound, err := s.adminRepo.GetById(ctx, *command.AdminId)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	if adminFound == nil {
		return nil, response2.NewDataNotFoundError("admin not found")
	}

	// 1. Update admin status or role
	updateAdminEntity := &adminEntity.AdminUpdate{
		Status: command.Status,
		Role:   command.Role,
	}

	if err = updateAdminEntity.ValidateAdminUpdate(); err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	adminEntity, err := s.adminRepo.UpdateOne(ctx, *command.AdminId, updateAdminEntity)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	// 2. Map to result
	return &adminCommand.UpdateAdminForSuperAdminCommandResult{
		Admin: mapper.NewAdminResult(adminEntity),
	}, nil
}

func (s *sSuperAdmin) GetOneAdmin(
	ctx context.Context,
	query *adminQuery.GetOneAdminQuery,
) (result *adminQuery.AdminQueryResult, err error) {
	// 1. Get admin info
	adminFound, err := s.adminRepo.GetById(ctx, query.AdminId)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	if adminFound == nil {
		return nil, response2.NewDataNotFoundError("admin not found")
	}

	// 2. Map to result
	return &adminQuery.AdminQueryResult{
		Admin: mapper.NewAdminResult(adminFound),
	}, nil
}

func (s *sSuperAdmin) GetManyAdmin(
	ctx context.Context,
	query *adminQuery.GetManyAdminQuery,
) (result *adminQuery.AdminQueryListResult, err error) {
	// 1. Get list admin
	adminEntities, paging, err := s.adminRepo.GetMany(ctx, query)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	var adminResults []*common.AdminResult
	for _, admin := range adminEntities {
		adminResult := mapper.NewAdminResult(admin)
		adminResults = append(adminResults, adminResult)
	}

	return &adminQuery.AdminQueryListResult{
		Admins:         adminResults,
		PagingResponse: paging,
	}, nil
}
