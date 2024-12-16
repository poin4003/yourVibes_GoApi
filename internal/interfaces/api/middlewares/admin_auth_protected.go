package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
	"strings"
)

func AdminAuthProtected() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 1. Check authHeader
		if authHeader == "" {
			response.ErrorResponse(c, response.ErrInvalidToken, http.StatusUnauthorized, "Authorization header is empty")
			c.Abort()
			return
		}

		// 2. Take token from authHeader
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.ErrorResponse(c, response.ErrInvalidToken, http.StatusUnauthorized, "Authorization header is invalid")
			c.Abort()
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
			response.ErrorResponse(c, response.ErrInvalidToken, http.StatusForbidden, err.Error())
			c.Abort()
			return
		}

		// 4. Take claims from token
		adminId := token.Claims.(jwt.MapClaims)["id"]
		role := token.Claims.(jwt.MapClaims)["role"].(bool)

		c.Set("userId", adminId)
		c.Set("role", role)

		c.Next()
	}
}
