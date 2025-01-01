package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
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
	loginCommand *user_command.LoginCommand,
) (result *user_command.LoginCommandResult, err error) {
	result = &user_command.LoginCommandResult{}
	result.User = nil
	result.AccessToken = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Find User
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", loginCommand.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrCodeEmailOrPasswordIsWrong
			result.HttpStatusCode = http.StatusNotFound
			return result, err
		}
		return result, err
	}

	// 2. Return if account is blocked by admin
	if !userFound.Status {
		result.ResultCode = response.ErrCodeAccountBlockedByAdmin
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("this account has been blocked for violating our community standards")
	}

	// 3. Hash password
	if !crypto.CheckPasswordHash(loginCommand.Password, *userFound.Password) {
		result.ResultCode = response.ErrCodeEmailOrPasswordIsWrong
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("invalid credentials")
	}

	// 4. Put claims into token
	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	// 5. Generate token
	accessTokenGen, err := jwtutil.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtSecretKey)
	if err != nil {
		return result, fmt.Errorf("cannot create access token: %v", err)
	}

	// 5. Map to command result
	result.User = mapper.NewUserResultFromEntity(userFound)
	result.AccessToken = &accessTokenGen
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserAuth) Register(
	ctx context.Context,
	registerCommand *user_command.RegisterCommand,
) (result *user_command.RegisterCommandResult, err error) {
	result = &user_command.RegisterCommandResult{}
	// 1. Check user exist in user table
	userFound, err := s.userRepo.CheckUserExistByEmail(ctx, registerCommand.Email)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
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
	newUser, err := user_entity.NewUserLocal(
		registerCommand.FamilyName,
		registerCommand.Name,
		registerCommand.Email,
		hashedPassword,
		registerCommand.PhoneNumber,
		registerCommand.Birthday,
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

func (s *sUserAuth) AuthGoogle(
	ctx context.Context,
	command *user_command.AuthGoogleCommand,
) (result *user_command.AuthGoogleCommandResult, err error) {
	result = &user_command.AuthGoogleCommandResult{}
	result.User = nil
	result.AccessToken = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError

	// 1. Verify Google access token
	var googleTokenInfoUrl = global.Config.GoogleSetting.GoogleTokensInfoUrl
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?id_token=%s", googleTokenInfoUrl, command.OpenId), nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		result.ResultCode = response.ErrInvalidToken
		result.HttpStatusCode = http.StatusForbidden
		return result, fmt.Errorf("failed to verify openid: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.ResultCode = response.ErrInvalidToken
		result.HttpStatusCode = http.StatusForbidden
		return result, fmt.Errorf("invalid open id")
	}

	var tokenInfo common.TokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return result, fmt.Errorf("failed to decode token info: %w", err)
	}

	validClientIds := []string{
		global.Config.GoogleSetting.WebClientId,
		global.Config.GoogleSetting.AndroidClientId,
		global.Config.GoogleSetting.IosClientId,
	}

	clientIdValid := false
	for _, validID := range validClientIds {
		if tokenInfo.Aud == validID {
			clientIdValid = true
			break
		}
	}

	if !clientIdValid {
		result.ResultCode = response.ErrInvalidToken
		result.HttpStatusCode = http.StatusForbidden
		return result, fmt.Errorf("invalid client id")
	}

	// 3. Get user by email
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 2.1. Create new user
			newUser, err := user_entity.NewUserGoogle(
				command.FamilyName,
				command.Name,
				command.Email,
				command.AuthGoogleId,
				command.AvatarUrl,
			)
			if err != nil {
				return result, err
			}

			createdUser, err := s.userRepo.CreateOne(ctx, newUser)
			if err != nil {
				return result, err
			}

			// 2.2. Create setting for user
			newSetting, err := user_entity.NewSetting(createdUser.ID, consts.VI)
			if err != nil {
				return result, err
			}

			createdSetting, err := s.settingRepo.CreateOne(ctx, newSetting)
			if err != nil {
				return result, err
			}

			createdUser.Setting = createdSetting

			// 2.3. Validate user
			validatedUser, err := user_validator.NewValidatedUserForGoogleAuth(createdUser)
			if err != nil {
				return result, err
			}

			accessClaims := jwt.MapClaims{
				"id":  validatedUser.ID,
				"exp": time.Now().Add(time.Hour * 720).Unix(),
			}

			// 2.4. Generate token
			accessTokenGen, err := jwtutil.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtSecretKey)
			if err != nil {
				return result, fmt.Errorf("cannot create access token: %w", err)
			}
			result.User = mapper.NewUserResultFromValidateEntity(validatedUser)
			result.AccessToken = &accessTokenGen
			result.ResultCode = response.ErrCodeSuccess
			result.HttpStatusCode = http.StatusOK

			return result, nil
		}
		return result, err
	}

	// 3. Check authType
	// 3.1. Return if auth type is local
	if userFound.AuthType == consts.LOCAL_AUTH {
		result.ResultCode = response.ErrCodeInvalidLocalAuthType
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you can't use local account to google login")
	}

	// 4. Check auth google id
	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	// 4.1 Generate token
	accessTokenGen, err := jwtutil.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtSecretKey)
	if err != nil {
		return result, fmt.Errorf("cannot create access token: %w", err)
	}

	result.User = mapper.NewUserResultFromEntity(userFound)
	result.AccessToken = &accessTokenGen
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
