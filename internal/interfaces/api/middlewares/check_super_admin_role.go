package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

func CheckSuperAdminRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || role != true {
			ctx.Error(response.NewInvalidTokenError())
			return
		}

		ctx.Next()
	}
}
