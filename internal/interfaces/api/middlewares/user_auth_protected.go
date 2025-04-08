package middlewares

import (
	"fmt"
	"strings"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
)

type userAuthProtectedMiddleware struct {
	userInfoService services.IUserInfo
}

func NewUserAuthProtectedMiddleware(
	userInfoService services.IUserInfo,
) *userAuthProtectedMiddleware {
	return &userAuthProtectedMiddleware{
		userInfoService: userInfoService,
	}
}

type IUserAuthProtectedMiddleware interface {
	UserAuthProtected() gin.HandlerFunc
}

func (m *userAuthProtectedMiddleware) UserAuthProtected() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		// 1. Check authHeader
		if authHeader == "" {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		// 2. Take token from authHeader
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		tokenStr := tokenParts[1]
		secret := []byte(global.Config.Authentication.JwtSecretKey)

		// 3. Parse jwt and authenticate secret key
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		// 4. Take userId from token
		userIdStr, ok := token.Claims.(jwt.MapClaims)["id"].(string)
		if !ok {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		// 5. Check user by check userStatus service
		userStatus, err := m.userInfoService.GetUserStatusById(ctx, userId)
		if err != nil {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		// 6. Check user status
		if !*userStatus {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		ctx.Set("userId", userId)

		ctx.Next()
	}
}
