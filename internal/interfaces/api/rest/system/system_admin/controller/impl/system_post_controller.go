package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	postServices "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/system/system_admin/dto/request"
)

type cSystemAdminPost struct {
	newFeedService postServices.IPostNewFeed
}

func NewSystemAdminPostController(
	newFeedService postServices.IPostNewFeed,
) *cSystemAdminPost {
	return &cSystemAdminPost{
		newFeedService: newFeedService,
	}
}

// UpdatePostAndStatistics godoc
// @Summary system update like, comment, statistics count
// @Description system hack to update immediately like, comment, statistic of post to viral
// @Tags system_post
// @Accept json
// @Produce json
// @Param post_id path string true "PostId"
// @Security ApiKeyAuth
// @Router /systems/update_post_and_statistics/{post_id} [post]
func (c *cSystemAdminPost) UpdatePostAndStatistics(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	if err = c.newFeedService.UpdatePostAndStatistics(ctx, postId); err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	pkgResponse.OK(ctx, nil)
}

// DelayPostCreatedAt godoc
// @Summary system update created_at of post
// @Description system hack to update created_at of post delay to 7 days ago
// @Tags system_post
// @Accept json
// @Produce json
// @Param post_id path string true "PostId"
// @Security ApiKeyAuth
// @Router /systems/delay_post_created_at/{post_id} [post]
func (c *cSystemAdminPost) DelayPostCreatedAt(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	if err = c.newFeedService.DelayPostCreatedAt(ctx, postId); err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	pkgResponse.OK(ctx, nil)
}

// ExpiredAdvertiseByPostId godoc
// @Summary system expired advertise by post id
// @Description system hack to expired 1 of advertise by post id
// @Tags system_post
// @Accept json
// @Produce json
// @Param post_id path string true "PostId"
// @Security ApiKeyAuth
// @Router /systems/expired_advertise/{post_id} [post]
func (c *cSystemAdminPost) ExpiredAdvertiseByPostId(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	if err = c.newFeedService.ExpireAdvertiseByPostId(ctx, postId); err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	pkgResponse.OK(ctx, nil)
}

// PushAdvertiseToNewFeed godoc
// @Summary system push advertise to new feed
// @Description System push advertisement to new feed by numUsers
// @Tags system_post
// @Accept json
// @Produce json
// @Param input body request.NumUsers true "input"
// @Security ApiKeyAuth
// @Router /systems/push_advertise_to_new_feed [post]
func (c *cSystemAdminPost) PushAdvertiseToNewFeed(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	numUsers, ok := body.(*request.NumUsers)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid request type"))
		return
	}

	// 2. Call service
	err := c.newFeedService.PushAdvertisementToNewFeed(ctx, numUsers.NumUsers)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}

// PushFeaturePostToNewFeed godoc
// @Summary system push feature post to new feed
// @Description System push feature post to new feed by numUsers
// @Tags system_post
// @Accept json
// @Produce json
// @Param input body request.NumUsers true "input"
// @Security ApiKeyAuth
// @Router /systems/push_feature_post_to_new_feed [post]
func (c *cSystemAdminPost) PushFeaturePostToNewFeed(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	numUsers, ok := body.(*request.NumUsers)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid request type"))
		return
	}

	// 2. Call service
	err := c.newFeedService.PushFeaturePostToNewFeed(ctx, numUsers.NumUsers)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}

// CheckExpiryOfAdvertisement godoc
// @Summary system push check expiry of advertisement
// @Description System check and delete expiry advertisement from new feed
// @Tags system_post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /systems/check_expiry_of_advertisement [post]
func (c *cSystemAdminPost) CheckExpiryOfAdvertisement(ctx *gin.Context) {
	if err := c.newFeedService.CheckExpiryOfAdvertisement(ctx); err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}

// CheckExpiryOfFeaturePost godoc
// @Summary system push check expiry of feature post
// @Description System check and delete expiry feature post from new feed
// @Tags system_post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /systems/check_expiry_of_feature_post [post]
func (c *cSystemAdminPost) CheckExpiryOfFeaturePost(ctx *gin.Context) {
	if err := c.newFeedService.CheckExpiryOfFeaturePost(ctx); err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}
