package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
	"reflect"
)

func ValidateJsonBody(
	dto interface{},
	validateFunc func(interface{}) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dtoInstance := reflect.New(reflect.TypeOf(dto).Elem()).Interface()

		if err := ctx.ShouldBindJSON(dtoInstance); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		if err := validateFunc(dtoInstance); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("validatedRequest", dtoInstance)
		ctx.Next()
	}
}

func ValidateFormBody(
	dto interface{},
	validateFunc func(interface{}) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dtoInstance := reflect.New(reflect.TypeOf(dto).Elem()).Interface()

		if err := ctx.ShouldBind(dtoInstance); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		if err := validateFunc(dtoInstance); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("validatedRequest", dtoInstance)
		ctx.Next()
	}
}

func ValidateQuery(
	query interface{},
	validateFunc func(interface{}) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queryInstance := reflect.New(reflect.TypeOf(query).Elem()).Interface()

		if err := ctx.ShouldBindQuery(queryInstance); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		if err := validateFunc(queryInstance); err != nil {
			response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("validatedQuery", queryInstance)
		ctx.Next()
	}
}
