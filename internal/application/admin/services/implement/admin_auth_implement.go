package implement

import (
	"context"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/crypto"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/jwtutil"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	adminCommand "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	adminEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	adminRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sAdminAuth struct {
	adminRepo adminRepo.IAdminRepository
}

func NewAdminAuthImplement(
	adminRepo adminRepo.IAdminRepository,
) *sAdminAuth {
	return &sAdminAuth{
		adminRepo: adminRepo,
	}
}

func (s *sAdminAuth) Login(
	ctx context.Context,
	command *adminCommand.LoginCommand,
) (result *adminCommand.LoginCommandResult, err error) {
	// 1. Find admin
	adminFound, err := s.adminRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	if adminFound == nil {
		return nil, response2.NewDataNotFoundError("admin not found")
	}

	// 2. Return if account is blocked by admin
	if !adminFound.Status {
		return nil, response2.NewCustomError(response2.ErrCodeAccountBlockedBySuperAdmin)
	}

	// 3. Check hash password
	if !crypto.CheckPasswordHash(command.Password, adminFound.Password) {
		return nil, response2.NewCustomError(response2.ErrCodeEmailOrPasswordIsWrong)
	}

	// 4. Put claims into token
	accessClaims := jwt.MapClaims{
		"id":   adminFound.ID,
		"role": adminFound.Role,
		"exp":  time.Now().Add(time.Hour * 720).Unix(),
	}

	// 5. Generate token
	accessTokenGen, err := jwtutil.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtAdminSecretKey)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	// 6. Map to result
	return &adminCommand.LoginCommandResult{
		Admin:       mapper.NewAdminResult(adminFound),
		AccessToken: &accessTokenGen,
	}, nil
}

func (s *sAdminAuth) ChangeAdminPassword(
	ctx context.Context,
	command *adminCommand.ChangeAdminPasswordCommand,
) error {
	// 1. Find admin
	adminFound, err := s.adminRepo.GetById(ctx, command.AdminId)
	if err != nil {
		return response2.NewServerFailedError(err.Error())
	}

	if adminFound == nil {
		return response2.NewDataNotFoundError("admin not found")
	}

	// 2. Check old password
	if !crypto.CheckPasswordHash(command.OldPassword, adminFound.Password) {
		return response2.NewCustomError(response2.ErrCodeOldPasswordIsWrong)
	}

	// 3. Update new password
	hashedPassword, err := crypto.HashPassword(command.NewPassword)
	if err != nil {
		return response2.NewServerFailedError(err.Error())
	}

	updateAdminData := &adminEntity.AdminUpdate{
		Password: pointer.Ptr(hashedPassword),
	}
	if err := updateAdminData.ValidateAdminUpdate(); err != nil {
		return response2.NewServerFailedError(err.Error())
	}

	_, err = s.adminRepo.UpdateOne(ctx, command.AdminId, updateAdminData)
	if err != nil {
		return response2.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sAdminAuth) ForgotAdminPassword(
	ctx context.Context,
	command *adminCommand.ForgotAdminPasswordCommand,
) error {
	// 1. Check admin exist
	adminFound, err := s.adminRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		return response2.NewServerFailedError(err.Error())
	}

	if adminFound == nil {
		return response2.NewDataNotFoundError("admin not found")
	}

	// 2. Update new password
	hashedPassword, err := crypto.HashPassword(command.NewPassword)
	if err != nil {
		return response2.NewServerFailedError(err.Error())
	}

	updateAdminData := &adminEntity.AdminUpdate{
		Password: pointer.Ptr(hashedPassword),
	}

	if err = updateAdminData.ValidateAdminUpdate(); err != nil {
		return response2.NewServerFailedError(err.Error())
	}

	_, err = s.adminRepo.UpdateOne(ctx, adminFound.ID, updateAdminData)
	if err != nil {
		return response2.NewServerFailedError(err.Error())
	}

	return nil
}
