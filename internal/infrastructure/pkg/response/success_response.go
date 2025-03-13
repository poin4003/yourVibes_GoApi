package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseDataWithPaging struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Paging  PagingResponse `json:"paging,omitempty"`
}

type PagingResponse struct {
	Limit int   `json:"limit"`
	Page  int   `json:"page"`
	Total int64 `json:"total"`
}

func sendSuccessResponse(
	ctx *gin.Context,
	code int,
	httpStatusCode int,
	message string,
	data interface{},
) {
	ctx.JSON(httpStatusCode, ResponseData{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func sendSuccessPagingResponse(
	ctx *gin.Context,
	code int,
	httpStatusCode int,
	message string,
	data interface{},
	paging PagingResponse,
) {
	ctx.JSON(httpStatusCode, ResponseDataWithPaging{
		Code:    code,
		Message: message,
		Data:    data,
		Paging:  paging,
	})
}

func SuccessResponse(ctx *gin.Context, code int, data interface{}) {
	if customCode, exists := GetCustomCode(code); exists {
		sendSuccessResponse(
			ctx, customCode.Code, customCode.HttpStatusCode, customCode.Message, data,
		)
	} else {
		sendSuccessResponse(
			ctx, customCode.Code, http.StatusOK, "Success", data,
		)
	}
}

func SuccessPagingResponse(ctx *gin.Context, code int, data interface{}, paging PagingResponse) {
	if customCode, exists := GetCustomCode(code); exists {
		sendSuccessPagingResponse(
			ctx, customCode.Code, customCode.HttpStatusCode, customCode.Message, data, paging,
		)
	} else {
		sendSuccessPagingResponse(
			ctx, customCode.Code, http.StatusOK, "Success", data, paging,
		)
	}
}

func OK(ctx *gin.Context, data interface{}) {
	SuccessResponse(ctx, ErrCodeSuccess, data)
}

func OKWithPaging(ctx *gin.Context, data interface{}, paging PagingResponse) {
	SuccessPagingResponse(ctx, ErrCodeSuccess, data, paging)
}
