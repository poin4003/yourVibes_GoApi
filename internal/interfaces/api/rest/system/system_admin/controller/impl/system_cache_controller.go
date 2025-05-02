package impl

import (
	"github.com/gin-gonic/gin"
	commentServices "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	postServices "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	userServices "github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type cSystemAdminCache struct {
	userService    userServices.IUserInfo
	postService    postServices.IPostUser
	commentService commentServices.ICommentUser
}

func NewSystemAdminCacheController(
	userService userServices.IUserInfo,
	postService postServices.IPostUser,
	commentService commentServices.ICommentUser,
) *cSystemAdminCache {
	return &cSystemAdminCache{
		userService:    userService,
		postService:    postService,
		commentService: commentService,
	}
}

// ClearAllCache godoc
// @Summary system cache clear
// @Description When admin need to clear all system cache
// @Tags system_cache
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /systems/all_cache [post]
func (c *cSystemAdminCache) ClearAllCache(ctx *gin.Context) {
	if err := c.userService.DeleteAllCache(ctx); err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	pkgResponse.OK(ctx, nil)
}

// ClearAllPostCache godoc
// @Summary system post cache clear
// @Description When admin need to clear all post system cache
// @Tags system_cache
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /systems/post_cache [post]
func (c *cSystemAdminCache) ClearAllPostCache(ctx *gin.Context) {
	if err := c.postService.ClearAllPostCaches(ctx); err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	pkgResponse.OK(ctx, nil)
}

// ClearAllCommentCache godoc
// @Summary system comment cache clear
// @Description When admin need to clear all comment system cache
// @Tags system_cache
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /systems/comment_cache [post]
func (c *cSystemAdminCache) ClearAllCommentCache(ctx *gin.Context) {
	if err := c.commentService.ClearAllCommentCaches(ctx); err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	pkgResponse.OK(ctx, nil)
}
