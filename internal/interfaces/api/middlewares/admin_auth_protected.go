package middlewares

import (
	"fmt"
	"strings"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
)

type adminAuthProtectedMiddleware struct {
	adminInfoService services.IAdminInfo
}

func NewAdminAuthProtectedMiddleware(
	adminInfoService services.IAdminInfo,
) *adminAuthProtectedMiddleware {
	return &adminAuthProtectedMiddleware{
		adminInfoService: adminInfoService,
	}
}

type IAdminAuthProtectedMiddleware interface {
	AdminAuthProtected() gin.HandlerFunc
}

func (m *adminAuthProtectedMiddleware) AdminAuthProtected() gin.HandlerFunc {
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
		secret := []byte(global.Config.Authentication.JwtAdminSecretKey)

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

		// 4. Take claims from token
		adminIdStr, ok := token.Claims.(jwt.MapClaims)["id"].(string)
		if !ok {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		adminId, err := uuid.Parse(adminIdStr)
		if err != nil {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		// 5. Check admin form db
		adminStatus, err := m.adminInfoService.GetAdminStatusById(ctx, adminId)
		if err != nil {
			ctx.Error(response.NewServerFailedError())
			ctx.Abort()
			return
		}

		// 6. Check admin status
		if !*adminStatus {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		roleClaim, ok := token.Claims.(jwt.MapClaims)["role"]
		if !ok {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		role, ok := roleClaim.(bool)
		if !ok {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		ctx.Set("adminId", adminId)
		ctx.Set("role", role)

		ctx.Next()
	}
}
