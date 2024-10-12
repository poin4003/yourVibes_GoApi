package response

import (
	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	Error ErrResponseChild `json:"error"`
}

type ResponseData struct {
	Code    int            `json:"code"`             // Status code
	Message string         `json:"message"`          // Status message
	Data    interface{}    `json:"data"`             // Data
	Paging  PagingResponse `json:"paging,omitempty"` // Paging (optional)
}

type ErrResponseChild struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	DetailErr string `json:"detail_err"`
}

type PagingResponse struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Total int64 `json:"total"`
}

func SuccessResponse(c *gin.Context, code int, httpStatus int, data interface{}) {
	c.JSON(httpStatus, ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}

func SuccessPagingResponse(c *gin.Context, code int, httpStatus int, data interface{}, paging PagingResponse) {
	c.JSON(httpStatus, ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
		Paging:  paging,
	})
}

func ErrorResponse(c *gin.Context, code int, httpStatus int, err string) {
	c.JSON(httpStatus, ErrResponse{
		Error: ErrResponseChild{
			Code:      code,
			Message:   msg[code],
			DetailErr: err,
		},
	})
}
