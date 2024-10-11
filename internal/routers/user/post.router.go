package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/post_user"
	"github.com/poin4003/yourVibes_GoApi/internal/middlewares/authentication"
)

type PostRouter struct{}

func (pr *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
	// Public router

	//userRouterPublic := Router.Group("/posts")
	//{
	//}

	// Private router
	postRouterPrivate := Router.Group("/posts")
	postRouterPrivate.Use(authentication.AuthProteced())
	{
		postRouterPrivate.POST("/createPost", post_user.PostUser.CreatePost)
		//postRouterPrivate.POST("/createPost", func(c *gin.Context) {
		//	c.JSON(200, gin.H{
		//		"status": "Ok",
		//	})
		//})
	}
}
