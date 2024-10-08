package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"log"
	"strings"
)

func AuthProtected() gin.HandlerFunc {
	db := global.Pdb
	cf := global.Config

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			log.Warn("empty authorization header")
			ctx.JSON(401, gin.H{
				"status":  "fail",
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warn("invalid token parts")
			ctx.JSON(401, gin.H{
				"status":  "fail",
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		tokenStr := tokenParts[1]
		secret := cf.Authentication.JwtScretKey

		claims, err := utils.VerifyJWT(tokenStr, secret)
		if err != nil {
			log.Warnf("invalid token: %v", err)
			ctx.JSON(401, gin.H{
				"status":  "fail",
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		userIdFloat := claims["id"].(float64)
		userId := uint(userIdFloat)
		role := claims["role"]

		if err := db.Model(&model.User{}).Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("user not found in the db")
			ctx.JSON(401, gin.H{
				"status":  "fail",
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		// Set user ID and role in context
		ctx.Set("userId", userId)
		ctx.Set("role", role)
		ctx.Next()
	}
}
