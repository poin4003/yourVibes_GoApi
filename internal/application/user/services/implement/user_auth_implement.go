package implement

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/third_party_authentication"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	userValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/validator"
	userRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
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
	userRepo    userRepo.IUserRepository
	settingRepo userRepo.ISettingRepository
}

func NewUserLoginImplement(
	userRepo userRepo.IUserRepository,
	settingRepo userRepo.ISettingRepository,
) *sUserAuth {
	return &sUserAuth{
		userRepo:    userRepo,
		settingRepo: settingRepo,
	}
}

func (s *sUserAuth) Login(
	ctx context.Context,
	loginCommand *userCommand.LoginCommand,
) (result *userCommand.LoginCommandResult, err error) {
	result = &userCommand.LoginCommandResult{
		User:           nil,
		AccessToken:    nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Find User
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", loginCommand.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrCodeEmailOrPasswordIsWrong
			result.HttpStatusCode = http.StatusBadRequest
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

	// 3. Check auth type
	if userFound.AuthType != consts.LOCAL_AUTH {
		result.ResultCode = response.ErrCodeEmailOrPasswordIsWrong
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("invalid auth type")
	}

	// 4. Hash password
	if !crypto.CheckPasswordHash(loginCommand.Password, *userFound.Password) {
		result.ResultCode = response.ErrCodeEmailOrPasswordIsWrong
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("invalid credentials")
	}

	// 5. Put claims into token
	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	// 6. Generate token
	accessTokenGen, err := jwtutil.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtSecretKey)
	if err != nil {
		return result, fmt.Errorf("cannot create access token: %v", err)
	}

	// 7. Map to command result
	result.User = mapper.NewUserResultFromEntity(userFound)
	result.AccessToken = &accessTokenGen
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserAuth) Register(
	ctx context.Context,
	command *userCommand.RegisterCommand,
) (result *userCommand.RegisterCommandResult, err error) {
	result = &userCommand.RegisterCommandResult{
		ResultCode: response.ErrServerFailed,
	}
	// 1. Check user exist in user table
	userFound, err := s.userRepo.CheckUserExistByEmail(ctx, command.Email)
	if err != nil {
		return result, err
	}

	if userFound {
		result.ResultCode = response.ErrCodeUserHasExists
		return result, fmt.Errorf("user %s already exists", command.Email)
	}

	// 3. Get Otp from Redis
	hashEmail := crypto.GetHash(strings.ToLower(command.Email))
	userKey := utils.GetUserKey(hashEmail)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			result.ResultCode = response.ErrCodeOtpNotExists
			return result, fmt.Errorf("no OTP found for %s", command.Email)
		}
		result.ResultCode = response.ErrCodeOtpNotExists
		return result, err
	}

	// 3. Compare Otp
	if otpFound != command.Otp {
		result.ResultCode = response.ErrInvalidOTP
		return result, fmt.Errorf("otp does not match for %s", command.Email)
	}

	// 4. Hash password
	hashedPassword, err := crypto.HashPassword(command.Password)
	if err != nil {
		return result, err
	}

	// 5. Create new user
	newUser, err := userEntity.NewUserLocal(
		command.FamilyName,
		command.Name,
		command.Email,
		hashedPassword,
		command.PhoneNumber,
		command.Birthday,
	)
	if err != nil {
		return result, err
	}

	createdUser, err := s.userRepo.CreateOne(ctx, newUser)
	if err != nil {
		return result, err
	}

	// 6. Create setting for user
	newSetting, err := userEntity.NewSetting(createdUser.ID, consts.VI)
	if err != nil {
		return result, err
	}

	createdSetting, err := s.settingRepo.CreateOne(ctx, newSetting)
	if err != nil {
		return result, err
	}

	createdUser.Setting = createdSetting

	// 7. Validate user
	validatedUser, err := userValidator.NewValidatedUser(createdUser)
	if err != nil {
		return result, err
	}

	// 8. Send success email for user
	if err = sendto.SendTemplateEmail(
		[]string{command.Email},
		consts.HOST_EMAIL,
		"sign_up_success.html",
		map[string]interface{}{"email": command.Email},
		"Yourvibes sign up successful",
	); err != nil {
		result.ResultCode = response.ErrSendEmailOTP
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
	case errors.Is(err, redis.Nil):
		fmt.Println("Key does not exist")
	case err != nil:
		fmt.Println("Get failed::", err)
		return response.ErrCodeOtpNotExists, err
	case otpFound != "":
		return response.ErrCodeOtpNotExists, fmt.Errorf("otp already exists but not registered")
	}

	// 4. Generate OTP
	otpNew := random.GenerateSixDigitOtp()

	// 5. save OTP into Redis with expiration time
	err = global.Rdb.SetEx(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	if err != nil {
		return response.ErrInvalidOTP, err
	}

	// 6. send OTP
	err = sendto.SendTemplateEmail(
		[]string{email},
		consts.HOST_EMAIL,
		"otp_auth.html",
		map[string]interface{}{"otp": strconv.Itoa(otpNew)},
		"Yourvibes OTP Verification",
	)

	if err != nil {
		return response.ErrSendEmailOTP, err
	}

	return response.ErrCodeSuccess, nil
}

func (s *sUserAuth) ChangePassword(
	ctx context.Context,
	command *userCommand.ChangePasswordCommand,
) (result *userCommand.ChangePasswordCommandResult, err error) {
	result = &userCommand.ChangePasswordCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Find user
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.UserNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check auth type
	if userFound.AuthType != consts.LOCAL_AUTH {
		result.ResultCode = response.ErrCodeInvalidLocalAuthType
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	// 3. Check old password
	if !crypto.CheckPasswordHash(command.OldPassword, *userFound.Password) {
		result.ResultCode = response.ErrCodeOldPasswordIsWrong
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("old password is wrong")
	}

	// 4. Update new password
	hashedPassword, err := crypto.HashPassword(command.NewPassword)
	if err != nil {
		return result, err
	}

	updateUserData := &userEntity.UserUpdate{
		Password: pointer.Ptr(hashedPassword),
	}
	if err := updateUserData.ValidateUserUpdate(); err != nil {
		return result, err
	}

	_, err = s.userRepo.UpdateOne(ctx, command.UserId, updateUserData)
	if err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserAuth) GetOtpForgotUserPassword(
	ctx context.Context,
	command *userCommand.GetOtpForgotUserPasswordCommand,
) (result *userCommand.GetOtpForgotUserPasswordCommandResult, err error) {
	result = &userCommand.GetOtpForgotUserPasswordCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Hash Email
	hashEmail := crypto.GetHash(strings.ToLower(command.Email))

	// 2. Check user exist
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.UserNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 3. Check auth type
	if userFound.AuthType != consts.LOCAL_AUTH {
		result.ResultCode = response.ErrCodeInvalidLocalAuthType
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you can't use forgot password if auth type is googleOAuth")
	}

	// 4. Check OTP exists
	userKey := utils.GetOtpForgotPasswordUser(hashEmail)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	switch {
	case errors.Is(err, redis.Nil):
		fmt.Println("Key does not exist")
	case err != nil:
		fmt.Println("Get failed::", err)
		return result, err
	case otpFound != "":
		result.ResultCode = response.ErrCodeOtpNotExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("otp already exists but not registered")
	}

	// 5. Generate OTP
	otpNew := random.GenerateSixDigitOtp()

	// 6. Save OTP into Redis with expiration time
	if err = global.Rdb.SetEx(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_FORGOT_USER_PASSWORD)*time.Minute).Err(); err != nil {
		return result, err
	}

	// 7. Send OTP
	if err = sendto.SendTemplateEmail(
		[]string{command.Email},
		consts.HOST_EMAIL,
		"otp_forgot_password.html",
		map[string]interface{}{"otp": strconv.Itoa(otpNew)},
		"Yourvibes OTP Verification",
	); err != nil {
		result.ResultCode = response.ErrSendEmailOTP
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserAuth) ForgotUserPassword(
	ctx context.Context,
	command *userCommand.ForgotUserPasswordCommand,
) (result *userCommand.ForgotUserPasswordCommandResult, err error) {
	result = &userCommand.ForgotUserPasswordCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Check user exist
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.UserNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("user %s doesn't exists", command.Email)
		}
		return result, err
	}

	// 2. Check auth type
	if userFound.AuthType != consts.LOCAL_AUTH {
		result.ResultCode = response.ErrCodeInvalidLocalAuthType
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you can't use forgot password if auth type is googleOAuth")
	}

	// 3. Get Otp from Redis
	hashEmail := crypto.GetHash(strings.ToLower(command.Email))
	userKey := utils.GetOtpForgotPasswordUser(hashEmail)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			result.ResultCode = response.ErrCodeOtpNotExists
			return result, fmt.Errorf("no otp found for %s", command.Email)
		}
		result.ResultCode = response.ErrCodeOtpNotExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	// 4. Compare otp
	if otpFound != command.Otp {
		result.ResultCode = response.ErrInvalidOTP
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("otp does not match for %s", command.Email)
	}

	// 5. Update new password
	hashedPassword, err := crypto.HashPassword(command.NewPassword)
	if err != nil {
		return result, err
	}

	updateUserData := &userEntity.UserUpdate{
		Password: pointer.Ptr(hashedPassword),
	}

	if err = updateUserData.ValidateUserUpdate(); err != nil {
		return result, err
	}

	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, updateUserData)
	if err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserAuth) AuthGoogle(
	ctx context.Context,
	command *userCommand.AuthGoogleCommand,
) (result *userCommand.AuthGoogleCommandResult, err error) {
	result = &userCommand.AuthGoogleCommandResult{
		User:           nil,
		AccessToken:    nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Call api google to get openid ODCI
	idToken, err := third_party_authentication.GetGoogleIDToken(command.AuthorizationCode, command.Platform, command.RedirectUrl)
	if err != nil {
		result.ResultCode = response.ErrCodeGoogleAuth
		result.HttpStatusCode = http.StatusForbidden
		return result, err
	}

	// 2. Get claims from openid
	claims, err := jwtutil.DecodeGoogleIDToken(idToken)
	if err != nil {
		return result, err
	}

	// 3. Get user by email
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", claims.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 2.1. Create new user
			newUser, err := userEntity.NewUserGoogle(
				claims.FamilyName,
				claims.GivenName,
				claims.Email,
				claims.Sub,
				claims.Picture,
			)
			if err != nil {
				return result, err
			}

			createdUser, err := s.userRepo.CreateOne(ctx, newUser)
			if err != nil {
				return result, err
			}

			// 2.2. Create setting for user
			newSetting, err := userEntity.NewSetting(createdUser.ID, consts.VI)
			if err != nil {
				return result, err
			}

			createdSetting, err := s.settingRepo.CreateOne(ctx, newSetting)
			if err != nil {
				return result, err
			}

			createdUser.Setting = createdSetting

			// 2.3. Validate user
			validatedUser, err := userValidator.NewValidatedUserForGoogleAuth(createdUser)
			if err != nil {
				return result, err
			}

			// 2.4. Send success email to user
			if err = sendto.SendTemplateEmail(
				[]string{claims.Email},
				consts.HOST_EMAIL,
				"sign_up_success.html",
				map[string]interface{}{"email": claims.Email},
				"Yourvibes sign up with Google account successfully",
			); err != nil {
				result.ResultCode = response.ErrSendEmailOTP
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

	// 3. Return if account is blocked (status = false)
	if !userFound.Status {
		result.ResultCode = response.ErrCodeAccountBlockedByAdmin
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("this account has been blocked for violating our community standards")
	}

	// 4. Return if auth type is local
	if userFound.AuthType == consts.LOCAL_AUTH {
		result.ResultCode = response.ErrCodeInvalidLocalAuthType
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you can't use local account to google login")
	}

	// 5. Check auth google id
	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	// 5.1 Generate token
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

func (s *sUserAuth) AppAuthGoogle(
	ctx context.Context,
	command *userCommand.AuthAppGoogleCommand,
) (result *userCommand.AuthGoogleCommandResult, err error) {
	result = &userCommand.AuthGoogleCommandResult{
		User:           nil,
		AccessToken:    nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}

	// 1. Get claims from openid
	claims, err := jwtutil.DecodeGoogleIDToken(command.OpenId)
	if err != nil {
		return result, err
	}

	// 2. Get user by email
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", claims.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 2.1. Create new user
			newUser, err := userEntity.NewUserGoogle(
				claims.FamilyName,
				claims.GivenName,
				claims.Email,
				claims.Sub,
				claims.Picture,
			)
			if err != nil {
				return result, err
			}

			createdUser, err := s.userRepo.CreateOne(ctx, newUser)
			if err != nil {
				return result, err
			}

			// 2.2. Create setting for user
			newSetting, err := userEntity.NewSetting(createdUser.ID, consts.VI)
			if err != nil {
				return result, err
			}

			createdSetting, err := s.settingRepo.CreateOne(ctx, newSetting)
			if err != nil {
				return result, err
			}

			createdUser.Setting = createdSetting

			// 2.3. Validate user
			validatedUser, err := userValidator.NewValidatedUserForGoogleAuth(createdUser)
			if err != nil {
				return result, err
			}

			// 2.4. Send success email to user
			if err = sendto.SendTemplateEmail(
				[]string{claims.Email},
				consts.HOST_EMAIL,
				"sign_up_success.html",
				map[string]interface{}{"email": claims.Email},
				"Yourvibes sign up with Google Account successfully",
			); err != nil {
				result.ResultCode = response.ErrSendEmailOTP
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

	// 3. Return if account is blocked (status = false)
	if !userFound.Status {
		result.ResultCode = response.ErrCodeAccountBlockedByAdmin
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("this account has been blocked for violating our community standards")
	}

	// 4. Return if auth type is local
	if userFound.AuthType == consts.LOCAL_AUTH {
		result.ResultCode = response.ErrCodeInvalidLocalAuthType
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you can't use local account to google login")
	}

	// 5. Check auth google id
	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	// 6 Generate token
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
