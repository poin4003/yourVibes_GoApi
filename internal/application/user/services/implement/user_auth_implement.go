package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	user_validator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/validator"
	user_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/crypto"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/jwtutil"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/random"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/sendto"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type sUserAuth struct {
	userRepo    user_repo.IUserRepository
	settingRepo user_repo.ISettingRepository
}

func NewUserLoginImplement(
	userRepo user_repo.IUserRepository,
	settingRepo user_repo.ISettingRepository,
) *sUserAuth {
	return &sUserAuth{
		userRepo:    userRepo,
		settingRepo: settingRepo,
	}
}

func (s *sUserAuth) Login(
	ctx context.Context,
	loginCommand *command.LoginCommand,
) (result *command.LoginCommandResult, err error) {
	result = &command.LoginCommandResult{}
	// 1. Find User
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", loginCommand.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, err
	}

	// 2. Hash password
	if !crypto.CheckPasswordHash(loginCommand.Password, userFound.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 3. Put claims into token
	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	// 4. Generate token
	accessTokenGen, err := jwtutil.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtScretKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create access token: %v", err)
	}

	// 5. Map to command result
	result.User = mapper.NewUserResultFromEntity(userFound)
	result.AccessToken = accessTokenGen

	return result, nil
}

func (s *sUserAuth) Register(
	ctx context.Context,
	registerCommand *command.RegisterCommand,
) (result *command.RegisterCommandResult, err error) {
	result = &command.RegisterCommandResult{}
	// 1. Check user exist in user table
	userFound, err := s.userRepo.CheckUserExistByEmail(ctx, registerCommand.Email)
	if err != nil {
		result.ResultCode = response.ErrCodeUserHasExists
		return result, err
	}

	if userFound {
		result.ResultCode = response.ErrCodeUserHasExists
		return result, fmt.Errorf("user %s already exists", registerCommand.Email)
	}

	// 3. Get Otp from Redis
	hashEmail := crypto.GetHash(strings.ToLower(registerCommand.Email))
	userKey := utils.GetUserKey(hashEmail)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	if err != nil {
		if err == redis.Nil {
			result.ResultCode = response.ErrCodeOtpNotExists
			return result, fmt.Errorf("no OTP found for %s", registerCommand.Email)
		}
		result.ResultCode = response.ErrCodeOtpNotExists
		return result, err
	}

	// 3. Compare Otp
	if otpFound != registerCommand.Otp {
		result.ResultCode = response.ErrInvalidOTP
		return result, fmt.Errorf("otp does not match for %s", registerCommand.Email)
	}

	// 4. Hash password
	hashedPassword, err := crypto.HashPassword(registerCommand.Password)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		return result, err
	}

	// 5. Create new user
	newUser, err := user_entity.NewUser(
		registerCommand.FamilyName,
		registerCommand.Name,
		registerCommand.Email,
		hashedPassword,
		registerCommand.PhoneNumber,
		registerCommand.Birthday,
		consts.LOCAL_AUTH,
	)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		return result, err
	}

	createdUser, err := s.userRepo.CreateOne(ctx, newUser)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		return result, err
	}

	// 6. Create setting for user
	newSetting, err := user_entity.NewSetting(createdUser.ID, consts.VI)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		return result, err
	}

	createdSetting, err := s.settingRepo.CreateOne(ctx, newSetting)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		return result, err
	}

	createdUser.Setting = createdSetting

	// 7. Validate user
	validatedUser, err := user_validator.NewValidatedUser(createdUser)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		return result, err
	}

	result.User = mapper.NewUserResultFromValidateEntity(validatedUser)
	result.ResultCode = response.ErrCodeSuccess
	return result, nil
}

func (s *sUserAuth) VerifyEmail(
	ctx context.Context,
	email string,
) (resultCode int, err error) {
	// 1. hash Email
	hashEmail := crypto.GetHash(strings.ToLower(email))

	// 2. check user exists in users table
	userFound, err := s.userRepo.CheckUserExistByEmail(ctx, email)
	if err != nil {
		return response.ErrCodeUserHasExists, err
	}

	if userFound {
		return response.ErrCodeUserHasExists, fmt.Errorf("user %s already exists", email)
	}

	// 3. Check OTP exists
	userKey := utils.GetUserKey(hashEmail)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	switch {
	case err == redis.Nil:
		fmt.Println("Key does not exist")
	case err != nil:
		fmt.Println("Get failed::", err)
		return response.ErrCodeOtpNotExists, err
	case otpFound != "":
		return response.ErrCodeOtpNotExists, fmt.Errorf("otp %s already exists but not registered", otpFound)
	}

	// 4. Generate OTP
	otpNew := random.GenerateSixDigitOtp()

	// 5. save OTP into Redis with expiration time
	err = global.Rdb.SetEx(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	if err != nil {
		return response.ErrInvalidOTP, err
	}

	// 6. send OTP
	err = sendto.SendTemplateEmailOtp(
		[]string{email},
		consts.HOST_EMAIL,
		"otp-auth.html",
		map[string]interface{}{"otp": strconv.Itoa(otpNew)},
	)

	if err != nil {
		return response.ErrSendEmailOTP, err
	}

	return response.ErrCodeSuccess, nil
}
