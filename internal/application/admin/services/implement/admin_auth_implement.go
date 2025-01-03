package implement

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	adminCommand "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	adminEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	adminRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/crypto"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/jwtutil"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
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
	result = &adminCommand.LoginCommandResult{
		Admin:          nil,
		AccessToken:    nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Find admin
	adminFound, err := s.adminRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrCodeEmailOrPasswordIsWrong
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Return if account is blocked by admin
	if !adminFound.Status {
		result.ResultCode = response.ErrCodeAccountBlockedBySuperAdmin
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("this account has been blocked by super admin")
	}

	// 3. Check hash password
	if !crypto.CheckPasswordHash(command.Password, adminFound.Password) {
		result.ResultCode = response.ErrCodeEmailOrPasswordIsWrong
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("invalid credentials")
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
		return result, err
	}

	// 6. Map to result
	result.AccessToken = &accessTokenGen
	result.Admin = mapper.NewAdminResult(adminFound)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sAdminAuth) ChangeAdminPassword(
	ctx context.Context,
	command *adminCommand.ChangeAdminPasswordCommand,
) (result *adminCommand.ChangeAdminPasswordCommandResult, err error) {
	result = &adminCommand.ChangeAdminPasswordCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Find admin
	adminFound, err := s.adminRepo.GetById(ctx, command.AdminId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check old password
	if !crypto.CheckPasswordHash(command.OldPassword, adminFound.Password) {
		result.ResultCode = response.ErrCodeOldPasswordIsWrong
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("old password is wrong")
	}

	// 3. Update new password
	hashedPassword, err := crypto.HashPassword(command.NewPassword)
	if err != nil {
		return result, err
	}

	updateAdminData := &adminEntity.AdminUpdate{
		Password: pointer.Ptr(hashedPassword),
	}
	if err := updateAdminData.ValidateAdminUpdate(); err != nil {
		return result, err
	}

	_, err = s.adminRepo.UpdateOne(ctx, command.AdminId, updateAdminData)
	if err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sAdminAuth) ForgotAdminPassword(
	ctx context.Context,
	command *adminCommand.ForgotAdminPasswordCommand,
) (result *adminCommand.ForgotAdminPasswordCommandResult, err error) {
	result = &adminCommand.ForgotAdminPasswordCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Check admin exist
	adminFound, err := s.adminRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("admin %s doesn't exists", command.Email)
		}
		return result, err
	}

	// 2. Update new password
	hashedPassword, err := crypto.HashPassword(command.NewPassword)
	if err != nil {
		return result, err
	}

	updateAdminData := &adminEntity.AdminUpdate{
		Password: pointer.Ptr(hashedPassword),
	}

	if err = updateAdminData.ValidateAdminUpdate(); err != nil {
		return result, err
	}

	_, err = s.adminRepo.UpdateOne(ctx, adminFound.ID, updateAdminData)
	if err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
