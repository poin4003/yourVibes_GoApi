package implement

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils"
	crypto2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/crypto"
	jwtutil2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/jwtutil"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/random"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/sendto"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/third_party_authentication"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	userValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/validator"
	repository "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sUserAuth struct {
	userRepo      repository.IUserRepository
	settingRepo   repository.ISettingRepository
	newFeedRepo   repository.INewFeedRepository
	userAuthCache cache.IUserAuthCache
}

func NewUserLoginImplement(
	userRepo repository.IUserRepository,
	settingRepo repository.ISettingRepository,
	newFeedRepo repository.INewFeedRepository,
	userAuthCache cache.IUserAuthCache,
) *sUserAuth {
	return &sUserAuth{
		userRepo:      userRepo,
		settingRepo:   settingRepo,
		newFeedRepo:   newFeedRepo,
		userAuthCache: userAuthCache,
	}
}

func (s *sUserAuth) Login(
	ctx context.Context,
	loginCommand *userCommand.LoginCommand,
) (result *userCommand.LoginCommandResult, err error) {
	// 1. Find User
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", loginCommand.Email)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return nil, response.NewCustomError(response.ErrCodeEmailOrPasswordIsWrong)
	}

	// 2. Return if account is blocked by admin
	if !userFound.Status {
		return nil, response.NewCustomError(response.ErrCodeAccountBlockedByAdmin)
	}

	// 3. Check auth type
	if userFound.AuthType != consts.LOCAL_AUTH {
		return nil, response.NewCustomError(response.ErrCodeInvalidLocalAuthType)
	}

	// 4. Hash password
	if !crypto2.CheckPasswordHash(loginCommand.Password, *userFound.Password) {
		return nil, response.NewCustomError(response.ErrCodeEmailOrPasswordIsWrong)
	}

	// 5. Put claims into token
	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	// 6. Generate token
	accessTokenGen, err := jwtutil2.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtSecretKey)
	if err != nil {
		return result, response.NewServerFailedError("can not create access token")
	}

	// 7. Map to command result
	return &userCommand.LoginCommandResult{
		User:        mapper.NewUserResultFromEntity(userFound),
		AccessToken: &accessTokenGen,
	}, nil
}

func (s *sUserAuth) Register(
	ctx context.Context,
	command *userCommand.RegisterCommand,
) (result *userCommand.RegisterCommandResult, err error) {
	// 1. Check user exist in user table
	userFound, err := s.userRepo.CheckUserExistByEmail(ctx, command.Email)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if userFound {
		return nil, response.NewCustomError(
			response.ErrCodeUserHasExists,
			fmt.Sprintf("user %s already exists", command.Email),
		)
	}

	// 2. Get Otp from Redis
	hashEmail := crypto2.GetHash(strings.ToLower(command.Email))
	userKey := utils.GetUserKey(hashEmail)
	otpFound, err := s.userAuthCache.GetOtp(ctx, userKey)
	if err != nil {
		return nil, err
	}

	// 3. Compare Otp
	if *otpFound != command.Otp {
		return nil, response.NewCustomError(
			response.ErrInvalidOTP,
			fmt.Sprintf("otp does not match for %s", command.Email),
		)
	}

	// 4. Hash password
	hashedPassword, err := crypto2.HashPassword(command.Password)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
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
		return nil, response.NewServerFailedError(err.Error())
	}

	createdUser, err := s.userRepo.CreateOne(ctx, newUser)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 6. Create setting for user
	newSetting, err := userEntity.NewSetting(createdUser.ID, consts.VI)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	createdSetting, err := s.settingRepo.CreateOne(ctx, newSetting)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	createdUser.Setting = createdSetting

	// 7. Validate user
	validatedUser, err := userValidator.NewValidatedUser(createdUser)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 8. Send success email for user
	if err = sendto.SendTemplateEmail(
		[]string{command.Email},
		consts.HOST_EMAIL,
		"sign_up_success.html",
		map[string]interface{}{"email": command.Email},
		"Yourvibes sign up successful",
	); err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 9. Push post into new feed for user
	go func(userId uuid.UUID) {
		if err = s.newFeedRepo.CreateAdvertisePostsForUser(ctx, userId); err != nil {
			global.Logger.Error("Failed to push advertise for new user")
		}
		if err = s.newFeedRepo.CreateFeaturedPostsForUser(ctx, userId); err != nil {
			global.Logger.Error("Failed to push feature post for new user")
		}
	}(validatedUser.ID)

	return &userCommand.RegisterCommandResult{
		User: mapper.NewUserResultFromValidateEntity(validatedUser),
	}, nil
}

func (s *sUserAuth) VerifyEmail(
	ctx context.Context,
	email string,
) (err error) {
	// 1. hash Email
	hashEmail := crypto2.GetHash(strings.ToLower(email))

	// 2. check user exists in users table
	userFound, err := s.userRepo.CheckUserExistByEmail(ctx, email)
	if err != nil {
		return response.NewCustomError(response.ErrCodeUserHasExists)
	}

	if userFound {
		return response.NewCustomError(
			response.ErrCodeUserHasExists,
			fmt.Sprintf("user %s already exists", email),
		)
	}

	// 3. Check OTP exists
	userKey := utils.GetUserKey(hashEmail)
	otpFound, err := s.userAuthCache.GetOtp(ctx, userKey)
	if err != nil {
		return err
	}

	if *otpFound != "" {
		return response.NewCustomError(
			response.ErrCodeOtpNotExists,
			"otp already exists but not registered",
		)
	}

	// 4. Generate OTP
	otpNew := random.GenerateSixDigitOtp()

	// 5. save OTP into Redis with expiration time
	err = s.userAuthCache.SetOtp(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Minute)
	if err != nil {
		return err
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
		return response.NewCustomError(response.ErrSendEmailOTP, err.Error())
	}

	return nil
}

func (s *sUserAuth) ChangePassword(
	ctx context.Context,
	command *userCommand.ChangePasswordCommand,
) (err error) {
	// 1. Find user
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return response.NewCustomError(response.UserNotFound)
	}

	// 2. Check auth type
	if userFound.AuthType != consts.LOCAL_AUTH {
		return response.NewCustomError(response.ErrCodeInvalidLocalAuthType)
	}

	// 3. Check old password
	if !crypto2.CheckPasswordHash(command.OldPassword, *userFound.Password) {
		return response.NewCustomError(response.ErrCodeOldPasswordIsWrong)
	}

	// 4. Update new password
	hashedPassword, err := crypto2.HashPassword(command.NewPassword)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	updateUserData := &userEntity.UserUpdate{
		Password: pointer.Ptr(hashedPassword),
	}
	if err = updateUserData.ValidateUserUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.userRepo.UpdateOne(ctx, command.UserId, updateUserData)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sUserAuth) GetOtpForgotUserPassword(
	ctx context.Context,
	command *userCommand.GetOtpForgotUserPasswordCommand,
) (err error) {
	// 1. Hash Email
	hashEmail := crypto2.GetHash(strings.ToLower(command.Email))

	// 2. Check user exist
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return response.NewCustomError(response.UserNotFound)
	}

	// 3. Check auth type
	if userFound.AuthType != consts.LOCAL_AUTH {
		return response.NewCustomError(response.ErrCodeInvalidLocalAuthType)
	}

	// 4. Check OTP exists
	userKey := utils.GetOtpForgotPasswordUser(hashEmail)
	_, err = s.userAuthCache.GetOtp(ctx, userKey)

	// 5. Generate OTP
	otpNew := random.GenerateSixDigitOtp()

	// 6. Save OTP into Redis with expiration time
	if err = s.userAuthCache.SetOtp(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_FORGOT_USER_PASSWORD)*time.Minute); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 7. Send OTP
	if err = sendto.SendTemplateEmail(
		[]string{command.Email},
		consts.HOST_EMAIL,
		"otp_forgot_password.html",
		map[string]interface{}{"otp": strconv.Itoa(otpNew)},
		"Yourvibes OTP Verification",
	); err != nil {
		return response.NewCustomError(response.ErrSendEmailOTP)
	}

	return nil
}

func (s *sUserAuth) ForgotUserPassword(
	ctx context.Context,
	command *userCommand.ForgotUserPasswordCommand,
) (err error) {
	// 1. Check user exist
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", command.Email)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return response.NewCustomError(response.UserNotFound)
	}

	// 2. Check auth type
	if userFound.AuthType != consts.LOCAL_AUTH {
		return response.NewCustomError(response.ErrCodeInvalidLocalAuthType)
	}

	// 3. Get Otp from Redis
	hashEmail := crypto2.GetHash(strings.ToLower(command.Email))
	userKey := utils.GetOtpForgotPasswordUser(hashEmail)
	otpFound, err := s.userAuthCache.GetOtp(ctx, userKey)
	if err != nil {
		return err
	}

	// 4. Compare otp
	if *otpFound != command.Otp {
		return response.NewCustomError(
			response.ErrInvalidOTP,
			fmt.Sprintf("otp does not match for %s", command.Email),
		)
	}

	// 5. Update new password
	hashedPassword, err := crypto2.HashPassword(command.NewPassword)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	updateUserData := &userEntity.UserUpdate{
		Password: pointer.Ptr(hashedPassword),
	}

	if err = updateUserData.ValidateUserUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, updateUserData)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sUserAuth) AuthGoogle(
	ctx context.Context,
	command *userCommand.AuthGoogleCommand,
) (result *userCommand.AuthGoogleCommandResult, err error) {
	// 1. Call api google to get openid ODCI
	idToken, err := third_party_authentication.GetGoogleIDToken(command.AuthorizationCode, command.Platform, command.RedirectUrl)
	if err != nil {
		return nil, response.NewCustomError(response.ErrCodeGoogleAuth)
	}

	// 2. Get claims from openid
	claims, err := jwtutil2.DecodeGoogleIDToken(idToken)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 3. Get user by email
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", claims.Email)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		// 2.1. Create new user
		var newUser *userEntity.User
		newUser, err = userEntity.NewUserGoogle(
			claims.FamilyName,
			claims.GivenName,
			claims.Email,
			claims.Sub,
			claims.Picture,
		)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		newUser, err = s.userRepo.CreateOne(ctx, newUser)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 2.2. Create setting for user
		var newSetting *userEntity.Setting
		newSetting, err = userEntity.NewSetting(newUser.ID, consts.VI)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		newSetting, err = s.settingRepo.CreateOne(ctx, newSetting)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		newUser.Setting = newSetting

		// 2.3. Validate user
		var validatedUser *userValidator.ValidatedUser
		validatedUser, err = userValidator.NewValidatedUserForGoogleAuth(newUser)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 2.4. Send success email to user
		if err = sendto.SendTemplateEmail(
			[]string{claims.Email},
			consts.HOST_EMAIL,
			"sign_up_success.html",
			map[string]interface{}{"email": claims.Email},
			"Yourvibes sign up with Google account successfully",
		); err != nil {
			return nil, response.NewCustomError(response.ErrSendEmailOTP)
		}

		accessClaims := jwt.MapClaims{
			"id":  validatedUser.ID,
			"exp": time.Now().Add(time.Hour * 720).Unix(),
		}

		// 2.4. Generate token
		var accessTokenGen string
		accessTokenGen, err = jwtutil2.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtSecretKey)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 9. Push post into new feed for user
		go func(userId uuid.UUID) {
			if err = s.newFeedRepo.CreateAdvertisePostsForUser(ctx, userId); err != nil {
				global.Logger.Error("Failed to push advertise for new user")
			}
			if err = s.newFeedRepo.CreateFeaturedPostsForUser(ctx, userId); err != nil {
				global.Logger.Error("Failed to push feature post for new user")
			}
		}(validatedUser.ID)

		return &userCommand.AuthGoogleCommandResult{
			User:        mapper.NewUserResultFromValidateEntity(validatedUser),
			AccessToken: &accessTokenGen,
		}, nil
	}

	// 3. Return if account is blocked (status = false)
	if !userFound.Status {
		return nil, response.NewCustomError(response.ErrCodeAccountBlockedByAdmin)
	}

	// 4. Return if auth type is local
	if userFound.AuthType == consts.LOCAL_AUTH {
		return nil, response.NewCustomError(response.ErrCodeInvalidLocalAuthType)
	}

	// 5. Check auth google id
	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	// 5.1 Generate token
	accessTokenGen, err := jwtutil2.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtSecretKey)
	if err != nil {
		return result, fmt.Errorf("cannot create access token: %w", err)
	}

	return &userCommand.AuthGoogleCommandResult{
		User:        mapper.NewUserResultFromEntity(userFound),
		AccessToken: &accessTokenGen,
	}, nil
}

func (s *sUserAuth) AppAuthGoogle(
	ctx context.Context,
	command *userCommand.AuthAppGoogleCommand,
) (result *userCommand.AuthGoogleCommandResult, err error) {
	// 1. Get claims from openid
	claims, err := jwtutil2.DecodeGoogleIDToken(command.OpenId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 2. Get user by email
	userFound, err := s.userRepo.GetOne(ctx, "email = ?", claims.Email)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		// 2.1. Create new user
		newUser, err := userEntity.NewUserGoogle(
			claims.FamilyName,
			claims.GivenName,
			claims.Email,
			claims.Sub,
			claims.Picture,
		)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		createdUser, err := s.userRepo.CreateOne(ctx, newUser)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 2.2. Create setting for user
		newSetting, err := userEntity.NewSetting(createdUser.ID, consts.VI)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		createdSetting, err := s.settingRepo.CreateOne(ctx, newSetting)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		createdUser.Setting = createdSetting

		// 2.3. Validate user
		validatedUser, err := userValidator.NewValidatedUserForGoogleAuth(createdUser)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 2.4. Send success email to user
		if err = sendto.SendTemplateEmail(
			[]string{claims.Email},
			consts.HOST_EMAIL,
			"sign_up_success.html",
			map[string]interface{}{"email": claims.Email},
			"Yourvibes sign up with Google Account successfully",
		); err != nil {
			return nil, response.NewCustomError(response.ErrSendEmailOTP)
		}

		accessClaims := jwt.MapClaims{
			"id":  validatedUser.ID,
			"exp": time.Now().Add(time.Hour * 720).Unix(),
		}

		// 2.4. Generate token
		accessTokenGen, err := jwtutil2.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtSecretKey)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 9. Push post into new feed for user
		go func(userId uuid.UUID) {
			if err = s.newFeedRepo.CreateAdvertisePostsForUser(ctx, userId); err != nil {
				global.Logger.Error("Failed to push advertise for new user")
			}
			if err = s.newFeedRepo.CreateFeaturedPostsForUser(ctx, userId); err != nil {
				global.Logger.Error("Failed to push feature post for new user")
			}
		}(validatedUser.ID)

		return &userCommand.AuthGoogleCommandResult{
			User:        mapper.NewUserResultFromValidateEntity(validatedUser),
			AccessToken: &accessTokenGen,
		}, nil
	}

	// 3. Return if account is blocked (status = false)
	if !userFound.Status {
		return nil, response.NewCustomError(response.ErrCodeAccountBlockedByAdmin)
	}

	// 4. Return if auth type is local
	if userFound.AuthType == consts.LOCAL_AUTH {
		return nil, response.NewCustomError(response.ErrCodeInvalidLocalAuthType)
	}

	// 5. Check auth google id
	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	// 6 Generate token
	accessTokenGen, err := jwtutil2.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtSecretKey)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return &userCommand.AuthGoogleCommandResult{
		User:        mapper.NewUserResultFromEntity(userFound),
		AccessToken: &accessTokenGen,
	}, nil
}
