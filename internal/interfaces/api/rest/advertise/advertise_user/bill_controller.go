package advertise_user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
	"strings"
)

type cBill struct {
}

func NewBillController() *cBill {
	return &cBill{}
}

func (c *cBill) ConfirmPayment(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to ConfirmPaymentRequest
	confirmPaymentRequest, ok := queryInput.(*request.ConfirmPaymentRequest)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Redirect to frontend
	if confirmPaymentRequest.ExtraData != "" {
		parts := strings.Split(confirmPaymentRequest.ExtraData, "<splitText>")

		if len(parts) >= 3 {
			redirectUrl := parts[2]
			ctx.Redirect(http.StatusFound, redirectUrl)
		}
	}

	// 4. Check resultCode to response for user if it failed
	fmt.Println(confirmPaymentRequest.ResultCode)
	if confirmPaymentRequest.ResultCode != "0" {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Momo confirm payment failed")
		return
	}

	// 5. Call service to confirm payment
	confirmPaymentCommand, err := confirmPaymentRequest.ToConfirmPaymentCommand()
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Failed to confirm payment command")
		return
	}

	result, err := services.Bill().ConfirmPayment(ctx, confirmPaymentCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, result)
}
