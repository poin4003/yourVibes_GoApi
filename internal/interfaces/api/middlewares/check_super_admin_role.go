package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

func CheckSuperAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != true {
			response.ErrorResponse(c, response.ErrInvalidToken, http.StatusForbidden, "You must be super admin to access this function")
			c.Abort()
			return
		}

		c.Next()
	}
}
