package service_implement

import (
	"context"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/internal/utils"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/crypto"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/random"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/sendto"
	"github.com/poin4003/yourVibes_GoApi/internal/vo"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"time"
)

type sUserAuth struct {
	repo repository.IUserRepository
}

func NewUserLoginImplement(repo repository.IUserRepository) *sUserAuth {
	return &sUserAuth{repo: repo}
}

func (s *sUserAuth) Login(ctx context.Context, in *vo.LoginCredentials) (string, *model.User, error) {
	return "", &model.User{}, nil
}

func (s *sUserAuth) Register(ctx context.Context, in *vo.RegisterCredentials) (int, error) {
	// 1. check user exist in user table
	userFound, err := s.repo.CheckUserExistByEmail(ctx, in.Email)
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
		Birthday:    in.Birthday,
	}

	_, err = s.repo.CreateOne(ctx, user)
	if err != nil {
		return response.ErrCreateUserFail, err
	}

	return response.ErrCodeSuccess, nil
}

func (s *sUserAuth) VerifyEmail(ctx context.Context, email string) (int, error) {
	// 1. hash Email
	hashEmail := crypto.GetHash(strings.ToLower(email))

	// 2. check user exists in users table
	userFound, err := s.repo.CheckUserExistByEmail(ctx, email)
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
