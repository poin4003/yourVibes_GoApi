package helpers

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

func ValidateJsonBody(
	dto interface{},
	validateFunc func(interface{}) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dtoInstance := reflect.New(reflect.TypeOf(dto).Elem()).Interface()

		if err := ctx.ShouldBindJSON(dtoInstance); err != nil {
			ctx.Error(response.NewCustomError(response.ErrCodeValidate, err.Error()))
			return
		}

		if err := validateFunc(dtoInstance); err != nil {
			ctx.Error(response.NewCustomError(response.ErrCodeValidate, err.Error()))
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
			ctx.Error(response.NewCustomError(response.ErrCodeValidate, err.Error()))
			return
		}

		if err := validateFunc(dtoInstance); err != nil {
			ctx.Error(response.NewCustomError(response.ErrCodeValidate, err.Error()))
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
			ctx.Error(response.NewCustomError(response.ErrCodeValidate, err.Error()))
			return
		}

		if err := validateFunc(queryInstance); err != nil {
			ctx.Error(response.NewCustomError(response.ErrCodeValidate, err.Error()))
			return
		}

		ctx.Set("validatedQuery", queryInstance)
		ctx.Next()
	}
}
