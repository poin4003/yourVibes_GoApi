package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Override error
type CustomError struct {
	CustomResponse
	MessageDetail string `json:"message_detail,omitempty"`
}

func (e CustomError) Error() string {
	if e.MessageDetail != "" {
		return fmt.Sprintf("Code: %d, Message: %s, Detail: %s", e.Code, e.Message, e.MessageDetail)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

func NewCustomError(code int, messageDetail ...string) error {
	customCode, exist := GetCustomCode(code)
	if !exist {
		customCode = CustomResponse{
			Code:           ErrServerFailed,
			Message:        "Internal Server Error",
			HttpStatusCode: http.StatusInternalServerError,
		}
	}

	// if error has a messageDetail
	var detail string
	if len(messageDetail) > 0 {
		detail = messageDetail[0]
	}

	return CustomError{
		CustomResponse: customCode,
		MessageDetail:  detail,
	}
}

// Init custom error method
func NewServerFailedError(messageDetail ...string) error {
	return NewCustomError(ErrServerFailed, messageDetail...)
}

func NewInvalidTokenError(messageDetail ...string) error {
	return NewCustomError(ErrInvalidToken, messageDetail...)
}

func NewValidateError(messageDetail ...string) error {
	return NewCustomError(ErrCodeValidate, messageDetail...)
}

func NewDataNotFoundError(messageDetail ...string) error {
	return NewCustomError(ErrDataNotFound, messageDetail...)
}

// Method to reponse error
func ErrorResponse(ctx *gin.Context, code int, messageDetail ...string) {
	customCode, exists := GetCustomCode(code)
	if !exists {
		customCode = CustomResponse{
			Code:           ErrServerFailed,
			Message:        "Internal Server Error",
			HttpStatusCode: http.StatusInternalServerError,
		}
	}

	var detail string
	if len(messageDetail) > 0 {
		detail = messageDetail[0]
	}

	sendErrorResponse(ctx, customCode.HttpStatusCode, customCode.Code, customCode.Message, detail)
}

func sendErrorResponse(
	ctx *gin.Context,
	httpStatusCode int,
	code int,
	message string,
	messageDetail string,
) {
	response := gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	}

	if messageDetail != "" {
		response["error"].(gin.H)["message_detail"] = messageDetail
	}

	ctx.JSON(httpStatusCode, response)
}
