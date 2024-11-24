package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user"
)

type AdvertiseRouter struct{}

func (ar *AdvertiseRouter) InitAdvertiseRouter(Router *gin.RouterGroup) {
	// Public router
	advertiseController := advertise_user.NewAdvertiseController()

	// Private router
	advertiseRouterPrivate := Router.Group("/advertise")
	advertiseRouterPrivate.Use(middlewares.AuthProteced())
	{
		advertiseRouterPrivate.POST("/", advertiseController.CreateAdvertise)
	}
}
