package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
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
// @Security ApiKeyAuth
// @Router /posts/new_feeds/{post_id}/ [delete]
func (c *cPostNewFeed) DeleteNewFeed(ctx *gin.Context) {
	// 1. Get post id from param path
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call service
	deleteNewFeedCommand := &postCommand.DeleteNewFeedCommand{PostId: postId, UserId: userIdClaim}

	result, err := services.PostNewFeed().DeleteNewFeed(ctx, deleteNewFeedCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkgResponse.SuccessResponse(ctx, result.ResultCode, http.StatusOK, nil)
}

// GetNewFeeds godoc
// @Summary Get a list of new feed
// @Description Get a list of new feed
// @Tags post_new_feed
// @Param limit query int false "limit on page"
// @Param page query int false "current page"
// @Security ApiKeyAuth
// @Router /posts/new_feeds/ [get]
func (c *cPostNewFeed) GetNewFeeds(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to userQueryObject
	newFeedQueryObject, ok := queryInput.(*query.NewFeedQueryObject)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call services
	getNewFeedQuery, err := newFeedQueryObject.ToGetNewFeedQuery(userIdClaim)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.PostNewFeed().GetNewFeeds(ctx, getNewFeedQuery)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	var postDtos []*response.PostWithLikedDto
	for _, postResult := range result.Posts {
		postDtos = append(postDtos, response.ToPostWithLikedDto(*postResult))
	}

	pkgResponse.SuccessPagingResponse(ctx, result.ResultCode, http.StatusOK, postDtos, *result.PagingResponse)
}
