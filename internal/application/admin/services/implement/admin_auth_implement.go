package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	admin_command "github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/mapper"
	admin_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/crypto"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/jwtutil"
	"gorm.io/gorm"
	"time"
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
	command *admin_command.LoginCommand,
) (result *admin_command.LoginCommandResult, err error) {
	result = &admin_command.LoginCommandResult{}
	// 1. Find admin
	adminFound, err := s.adminRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, err
	}

	// 2. Check hash password
	if !crypto.CheckPasswordHash(command.Password, adminFound.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 3. Put claims into token
	accessClaims := jwt.MapClaims{
		"id":   adminFound.ID,
		"role": adminFound.Role,
		"exp":  time.Now().Add(time.Hour * 720).Unix(),
	}

	// 4. Generate token
	accessTokenGen, err := jwtutil.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtAdminSecretKey)
	if err != nil {
		return nil, err
	}

	// 5. Map to result
	result.AccessToken = accessTokenGen
	result.Admin = mapper.NewAdminResult(adminFound)

	return result, nil
}
