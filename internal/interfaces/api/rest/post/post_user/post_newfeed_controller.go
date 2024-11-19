package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cPostNewFeed struct{}

func NewPostNewFeedController() *cPostNewFeed {
	return &cPostNewFeed{}
}

// DeleteNewFeed godoc
// @Summary DeleteNewFeeds
// @Description delete new feeds
// @Tags post_new_feed
// @Param post_id path string true "post_id you want to delete over your newfeed"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/new_feeds/{post_id}/ [delete]
func (c *cPostNewFeed) DeleteNewFeed(ctx *gin.Context) {
	// 1. Get post id from param path
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call service
	deleteNewFeedCommand := &post_command.DeleteNewFeedCommand{PostId: postId, UserId: userIdClaim}

	result, err := services.PostNewFeed().DeleteNewFeed(ctx, deleteNewFeedCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, nil)
}

// GetNewFeeds godoc
// @Summary Get a list of new feed
// @Description Get a list of new feed
// @Tags post_new_feed
// @Param limit query int false "limit on page"
// @Param page query int false "current page"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/new_feeds/ [get]
func (c *cPostNewFeed) GetNewFeeds(ctx *gin.Context) {
	// 1. Validate and get query object from query
	var query query.NewFeedQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call services
	getNewFeedQuery, err := query.ToGetNewFeedQuery(userIdClaim)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.PostNewFeed().GetNewFeeds(ctx, getNewFeedQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, http.StatusOK, result.Posts, *result.PagingResponse)
}
