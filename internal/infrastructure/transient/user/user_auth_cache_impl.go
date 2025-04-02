package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/redis/go-redis/v9"
)

type tUserAuth struct {
	client *redis.Client
}

func NewUserAuthCache(client *redis.Client) *tUserAuth {
	return &tUserAuth{client: client}
}

func (t *tUserAuth) SetOtp(
	ctx context.Context,
	userKey, otp string, ttl time.Duration,
) error {
	if err := t.client.Set(ctx, userKey, otp, ttl).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return response.NewCustomError(
				response.ErrCodeOtpNotExists,
				fmt.Sprintf("no OTP found for %s", userKey),
			)
		}
		return response.NewServerFailedError(err.Error())
	}
	return nil
}

func (t *tUserAuth) GetOtp(
	ctx context.Context,
	userKey string,
) (*string, error) {
	otpFound, err := t.client.Get(ctx, userKey).Result()
	switch {
	case errors.Is(err, redis.Nil):
		fmt.Println("key does not exist")
	case err != nil:
		global.Logger.Error(err.Error())
		return nil, response.NewServerFailedError(err.Error())
	}

	return &otpFound, nil
}
