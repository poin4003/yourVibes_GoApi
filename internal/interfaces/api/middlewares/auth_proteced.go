package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func AuthProteced() gin.HandlerFunc {
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
		secret := []byte(global.Config.Authentication.JwtScretKey)

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

		// 4. Take userId from token
		userId := token.Claims.(jwt.MapClaims)["id"]

		// 5. Check exist user in db
		if err := global.Pdb.Model(&models.User{}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.ErrorResponse(c, response.ErrInvalidToken, http.StatusForbidden, err.Error())
				c.Abort()
				return
			}
			response.ErrorResponse(c, response.ErrInvalidToken, http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}

		c.Set("userId", userId)

		c.Next()
	}
}
