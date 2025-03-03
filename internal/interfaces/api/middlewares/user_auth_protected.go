package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

func UserAuthProtected() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		// 1. Check authHeader
		if authHeader == "" {
			ctx.Error(response.NewInvalidTokenError())
			return
		}

		// 2. Take token from authHeader
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.Error(response.NewInvalidTokenError())
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
			return
		}

		// 4. Take userId from token
		userIdStr, ok := token.Claims.(jwt.MapClaims)["id"].(string)
		if !ok {
			ctx.Error(response.NewInvalidTokenError())
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			ctx.Error(response.NewInvalidTokenError())
			return
		}

		// 5. Check user from db
		userStatus, err := services.UserInfo().GetUserStatusById(ctx, userId)
		if err != nil {
			ctx.Error(response.NewInvalidTokenError())
			return
		}

		// 6. Check user status
		if !userStatus {
			ctx.Error(response.NewInvalidTokenError())
			return
		}

		ctx.Set("userId", userId)

		ctx.Next()
	}
}
