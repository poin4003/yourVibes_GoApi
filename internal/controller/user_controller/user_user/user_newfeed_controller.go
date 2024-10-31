package user_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cUserNewFeed struct{}

func NewUserNewFeedController() *cUserNewFeed {
	return &cUserNewFeed{}
}

// DeleteNewFeed godoc
// @Summary DeleteNewFeeds
// @Description delete new feeds
// @Tags user_new_feed
// @Param post_id path string true "post_id you want to delete over your newfeed"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/new_feeds/{post_id}/ [delete]
func (c *cUserNewFeed) DeleteNewFeed(ctx *gin.Context) {
	// 1. Get post id from param path
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call service
	resultCode, httpStatusCode, err := services.UserNewFeed().DeleteNewFeed(ctx, userIdClaim, postId)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, nil)
}

// GetNewFeeds godoc
// @Summary Get a list of new feed
// @Description Get a list of new feed
// @Tags user_new_feed
// @Param limit query int false "limit on page"
// @Param page query int false "current page"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/new_feeds/ [get]
func (c *cUserNewFeed) GetNewFeeds(ctx *gin.Context) {
	// 1. Validate and get query object from query
	var query query_object.NewFeedQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call services
	postDtos, paging, resultCode, httpStatusCode, err := services.UserNewFeed().GetNewFeeds(ctx, userIdClaim, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessPagingResponse(ctx, resultCode, http.StatusOK, postDtos, *paging)
}
