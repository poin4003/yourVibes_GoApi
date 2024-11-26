package advertise_user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
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
		return
	}

	// 2. Convert to ConfirmPaymentRequest
	confirmPaymentRequest, ok := queryInput.(*request.ConfirmPaymentRequest)
	if !ok {
		return
	}

	// 3. Setup to get redirect url
	var redirectUrl string
	if confirmPaymentRequest.ExtraData != "" {
		parts := strings.Split(confirmPaymentRequest.ExtraData, "<splitText>")

		if len(parts) >= 3 {
			redirectUrl = parts[2]
		}
	}

	// 4. Check resultCode to response for user if it failed
	fmt.Println(confirmPaymentRequest.ResultCode)
	if confirmPaymentRequest.ResultCode != "0" {
		ctx.Redirect(http.StatusMovedPermanently, redirectUrl)
	}

	// 5. Call service to confirm payment
	confirmPaymentCommand, err := confirmPaymentRequest.ToConfirmPaymentCommand()
	if err != nil {
		ctx.Redirect(http.StatusMovedPermanently, redirectUrl)
	}

	_, err = services.Bill().ConfirmPayment(ctx, confirmPaymentCommand)
	if err != nil {
		ctx.Redirect(http.StatusMovedPermanently, redirectUrl)
	}

	ctx.Redirect(http.StatusMovedPermanently, redirectUrl)
	return
}
