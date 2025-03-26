package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user"
	advertiseRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
	advertiseQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/query"
)

type AdvertiseRouter struct{}

func (ar *AdvertiseRouter) InitAdvertiseRouter(Router *gin.RouterGroup) {
	// Public router
	advertiseController := advertise_user.NewAdvertiseController()
	billController := advertise_user.NewBillController()

	billRouterPublic := Router.Group("/bill")
	{
		billRouterPublic.GET("/",
			helpers.ValidateQuery(&advertiseRequest.ConfirmPaymentRequest{}, advertiseRequest.ValidateConfirmPaymentRequest),
			billController.ConfirmPayment,
		)
	}

	// Private router
	advertiseRouterPrivate := Router.Group("/advertise")
	advertiseRouterPrivate.Use(middlewares.UserAuthProtected())
	{
		advertiseRouterPrivate.POST("/",
			helpers.ValidateJsonBody(&advertiseRequest.CreateAdvertiseRequest{}, advertiseRequest.ValidateCreateAdvertiseRequest),
			advertiseController.CreateAdvertise,
		)

		advertiseRouterPrivate.GET("/",
			helpers.ValidateQuery(&advertiseQuery.AdvertiseQueryObject{}, advertiseQuery.ValidateAdvertiseQueryObject),
			advertiseController.GetManyAdvertise,
		)

		advertiseRouterPrivate.GET("/statistic/:advertise_id", advertiseController.GetAdvertiseWithStatistic)
	}
}
