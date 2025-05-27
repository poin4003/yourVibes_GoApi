package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// If server panic return 500
				fmt.Println("Panic đây: ", rec)
				response.ErrorResponse(ctx, response.ErrServerFailed)
				ctx.Abort()
			}
		}()

		ctx.Next()

		// If context has a error
		if len(ctx.Errors) > 0 {
			lastErr := ctx.Errors.Last()
			if customErr, ok := lastErr.Err.(response.CustomError); ok {
				response.ErrorResponse(ctx, customErr.Code, customErr.MessageDetail)
			} else {
				response.ErrorResponse(ctx, response.ErrServerFailed, lastErr.Error())
			}
			ctx.Abort()
		}
	}
}
