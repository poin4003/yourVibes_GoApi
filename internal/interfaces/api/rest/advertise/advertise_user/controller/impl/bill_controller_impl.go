package impl

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
)

type cBill struct {
	billService services.IBill
}

func NewBillController(
	billService services.IBill,
) *cBill {
	return &cBill{
		billService: billService,
	}
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

	err = c.billService.ConfirmPayment(ctx, confirmPaymentCommand)
	if err != nil {
		ctx.Redirect(http.StatusMovedPermanently, redirectUrl)
	}

	if strings.HasPrefix(redirectUrl, "exp://") || strings.HasPrefix(redirectUrl, "myapp://") {
		htmlContent := fmt.Sprintf(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>Redirecting...</title>
				<script type="text/javascript">
					window.location = "%s"; // Custom scheme redirect
				</script>
			</head>
			<body>
				<p>If you are not redirected, <a href="%s">click here</a>.</p>
			</body>
			</html>
		`, redirectUrl, redirectUrl)

		ctx.Data(http.StatusOK, "text/html", []byte(htmlContent))
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, redirectUrl)
}
