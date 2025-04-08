package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/media/controller"
)

type mediaRouter struct {
	mediaController controller.IMediaController
}

func NewMediaRouter(
	mediaController controller.IMediaController,
) *mediaRouter {
	return &mediaRouter{mediaController: mediaController}
}

func (r *mediaRouter) InitMediaRouter(router *gin.RouterGroup) {
	// 1. Private router
	mediaRouterPrivate := router.Group("/media")
	{
		mediaRouterPrivate.GET("/:file_name", r.mediaController.GetMedia)
	}
}
