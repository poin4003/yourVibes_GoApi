package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

func CheckSuperAdminRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || role != true {
			ctx.Error(response.NewInvalidTokenError())
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
