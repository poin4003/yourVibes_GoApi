package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
	"strings"
)

func UserAuthProtected() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		// 1. Check authHeader
		if authHeader == "" {
			response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, "Authorization header is empty")
			ctx.Abort()
			return
		}

		// 2. Take token from authHeader
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, "Authorization header is invalid")
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
			response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, err.Error())
			ctx.Abort()
			return
		}

		// 4. Take userId from token
		userIdStr, ok := token.Claims.(jwt.MapClaims)["id"].(string)
		if !ok {
			response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, "Invalid token")
			ctx.Abort()
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, "Invalid token")
			ctx.Abort()
			return
		}

		// 5. Check user from db
		userStatus, err := services.UserInfo().GetUserStatusById(ctx, userId)
		if err != nil {
			response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, "Invalid token")
			ctx.Abort()
			return
		}

		// 6. Check user status
		if !userStatus {
			response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, "Invalid token")
			ctx.Abort()
			return
		}

		ctx.Set("userId", userId)

		ctx.Next()
	}
}
