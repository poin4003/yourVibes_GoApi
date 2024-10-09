package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int         `json:"code"`    // Status code
	Message string      `json:"message"` // Status message
	Data    interface{} `json:"data"`    // Data
}

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int) {
	c.JSON(http.StatusBadRequest, ErrResponse{
		Code:    code,
		Message: msg[code],
	})
}
