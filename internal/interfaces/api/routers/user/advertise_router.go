package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/controller"
	advertiseRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
	advertiseQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/query"
)

type advertiseRouter struct {
	advertiseController     controller.IAdvertiseController
	billController          controller.IBillController
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware
}

func NewAdvertiseRouter(
	advertiseController controller.IAdvertiseController,
	billController controller.IBillController,
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware,
) *advertiseRouter {
	return &advertiseRouter{
		advertiseController:     advertiseController,
		billController:          billController,
		userProtectedMiddleware: userProtectedMiddleware,
	}
}

func (r *advertiseRouter) InitAdvertiseRouter(Router *gin.RouterGroup) {
	// Public router
	billRouterPublic := Router.Group("/bill")
	{
		billRouterPublic.GET("/",
			helpers.ValidateQuery(&advertiseRequest.ConfirmPaymentRequest{}, advertiseRequest.ValidateConfirmPaymentRequest),
			r.billController.ConfirmPayment,
		)
	}

	// Private router
	advertiseRouterPrivate := Router.Group("/advertise")
	advertiseRouterPrivate.Use(r.userProtectedMiddleware.UserAuthProtected())
	{
		advertiseRouterPrivate.POST("/",
			helpers.ValidateJsonBody(&advertiseRequest.CreateAdvertiseRequest{}, advertiseRequest.ValidateCreateAdvertiseRequest),
			r.advertiseController.CreateAdvertise,
		)

		advertiseRouterPrivate.GET("/",
			helpers.ValidateQuery(&advertiseQuery.AdvertiseQueryObject{}, advertiseQuery.ValidateAdvertiseQueryObject),
			r.advertiseController.GetManyAdvertise,
		)

		advertiseRouterPrivate.GET("/statistic/:advertise_id", r.advertiseController.GetAdvertiseWithStatistic)
	}
}
