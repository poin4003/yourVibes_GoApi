package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

func ValidateJsonBody(
	dto interface{},
	validateFunc func(interface{}) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.ShouldBindJSON(dto); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		if err := validateFunc(dto); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("validatedRequest", dto)
		ctx.Next()
	}
}

func ValidateFormBody(
	dto interface{},
	validateFunc func(interface{}) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.ShouldBind(dto); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		if err := validateFunc(dto); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("validatedRequest", dto)
		ctx.Next()
	}
}

func ValidateQuery(
	query interface{},
	validateFunc func(interface{}) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.ShouldBindQuery(query); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		if err := validateFunc(query); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("validatedQuery", query)
		ctx.Next()
	}
}
