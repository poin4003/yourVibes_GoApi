package middlewares

import (
	"github.com/gin-gonic/gin"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// If server panic return 500
				pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed)
				ctx.Abort()
			}
		}()

		ctx.Next()

		// If context has a error
		if len(ctx.Errors) > 0 {
			for _, err := range ctx.Errors {
				if customErr, ok := err.Err.(pkgResponse.CustomError); ok {
					pkgResponse.ErrorResponse(ctx, customErr.Code)
				} else {
					pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, customErr.MessageDetail)
				}
			}
			ctx.Abort()
		}
	}
}
