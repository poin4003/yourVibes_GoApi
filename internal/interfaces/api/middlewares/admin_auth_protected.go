package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
	"strings"
)

func AdminAuthProtected() gin.HandlerFunc {
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
		secret := []byte(global.Config.Authentication.JwtAdminSecretKey)

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

		// 4. Take claims from token
		adminIdStr, ok := token.Claims.(jwt.MapClaims)["id"].(string)
		if !ok {
			response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, "Invalid token")
			ctx.Abort()
			return
		}

		adminId, err := uuid.Parse(adminIdStr)
		if err != nil {
			response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, "Invalid token")
			ctx.Abort()
			return
		}

		// 5. Check admin form db
		adminStatus, err := services.AdminInfo().GetAdminStatusById(ctx, adminId)
		if err != nil {
			response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, "Invalid token")
			ctx.Abort()
			return
		}

		// 6. Check admin status
		if !adminStatus {
			response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, "Invalid token")
			ctx.Abort()
			return
		}

		role := token.Claims.(jwt.MapClaims)["role"].(bool)

		ctx.Set("adminId", adminId)
		ctx.Set("role", role)

		ctx.Next()
	}
}
