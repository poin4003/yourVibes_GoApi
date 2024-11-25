package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user"
	bill_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
)

type AdvertiseRouter struct{}

func (ar *AdvertiseRouter) InitAdvertiseRouter(Router *gin.RouterGroup) {
	// Public router
	advertiseController := advertise_user.NewAdvertiseController()
	billController := advertise_user.NewBillController()

	billRouterPublic := Router.Group("/bill")
	{
		billRouterPublic.GET("/",
			helpers.ValidateQuery(&bill_request.ConfirmPaymentRequest{}, bill_request.ValidateConfirmPaymentRequest),
			billController.ConfirmPayment,
		)
	}

	// Private router
	advertiseRouterPrivate := Router.Group("/advertise")
	advertiseRouterPrivate.Use(middlewares.AuthProteced())
	{
		advertiseRouterPrivate.POST("/", advertiseController.CreateAdvertise)
	}
}
