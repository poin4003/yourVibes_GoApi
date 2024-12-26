package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/media"
)

type MediaRouter struct{}

func (mr *MediaRouter) InitMediaRouter(router *gin.RouterGroup) {
	// 1. Init controller
	mediaController := media.NewMediaController()

	// 2. Private router
	mediaRouterPrivate := router.Group("/media")
	{
		mediaRouterPrivate.GET("/:file_name", mediaController.GetMedia)
	}
}
