package service_implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/auth_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/internal/utils"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/crypto"
	jwtutil "github.com/poin4003/yourVibes_GoApi/internal/utils/jwtutil"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/random"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/sendto"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type sUserAuth struct {
	userRepo    repository.IUserRepository
	settingRepo repository.ISettingRepository
}

func NewUserLoginImplement(
	userRepo repository.IUserRepository,
	settingRepo repository.ISettingRepository,
) *sUserAuth {
	return &sUserAuth{
		userRepo:    userRepo,
		settingRepo: settingRepo,
	}
}

func (s *sUserAuth) Login(
	ctx context.Context,
	in *auth_dto.LoginCredentials,
) (accessToken string, user *model.User, err error) {
	userFound, err := s.userRepo.GetUser(ctx, "email = ?", in.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, fmt.Errorf("invalid credentials")
		}
		return "", nil, err
	}

	if !crypto.CheckPasswordHash(in.Password, userFound.Password) {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	accessClaims := jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
	}

	accessTokenGen, err := jwtutil.GenerateJWT(accessClaims, jwt.SigningMethodHS256, global.Config.Authentication.JwtScretKey)
	if err != nil {
		return "", nil, fmt.Errorf("Cannot create access token: %v", err)
	}

	return accessTokenGen, userFound, nil
}

func (s *sUserAuth) Register(
	ctx context.Context,
	in *auth_dto.RegisterCredentials,
) (resultCode int, err error) {
	// 1. check user exist in user table
	userFound, err := s.userRepo.CheckUserExistByEmail(ctx, in.Email)
	if err != nil {
		return response.ErrCodeUserHasExists, err
	}

	if userFound {
		return response.ErrCodeUserHasExists, fmt.Errorf("user %s already exists", in.Email)
	}

	// 2. Get Otp from Redis
	hashEmail := crypto.GetHash(strings.ToLower(in.Email))
	userKey := utils.GetUserKey(hashEmail)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	if err != nil {
		if err == redis.Nil {
			return response.ErrCodeOtpNotExists, fmt.Errorf("no OTP found for %s", in.Email)
		}
		return response.ErrCodeOtpNotExists, err
	}

	// 3. compare Otp
	if otpFound != in.Otp {
		return response.ErrInvalidOTP, fmt.Errorf("otp does not match for %s", in.Email)
	}

	// 4. hash password
	hashedPassword, err := crypto.HashPassword(in.Password)
	if err != nil {
		return response.ErrHashPasswordFail, err
	}

	// 5. create new user
	user := &model.User{
		FamilyName:  in.FamilyName,
		Name:        in.Name,
		Email:       in.Email,
		Password:    hashedPassword,
		PhoneNumber: in.PhoneNumber,
		AvatarUrl:   "https://res.cloudinary.com/dkf51e57t/image/upload/v1728899949/yourVibes/tyvild61lxom0gdkfbm6.jpg",
		Birthday:    in.Birthday,
	}

	newUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return response.ErrCreateUserFail, err
	}

	// 6. create setting for user
	setting := &model.Setting{
		UserId:   newUser.ID,
		Language: consts.VI,
	}

	_, err = s.settingRepo.CreateSetting(ctx, setting)
	if err != nil {
		return response.ErrServerFailed, err
	}

	return response.ErrCodeSuccess, nil
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
