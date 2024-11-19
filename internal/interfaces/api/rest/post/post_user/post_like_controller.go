package post_user

import (
	"github.com/redis/go-redis/v9"
)

type PostLikeController struct {
	redisClient *redis.Client
}

func NewPostLikeController(
	redisClient *redis.Client,
) *PostLikeController {
	return &PostLikeController{
		redisClient: redisClient,
	}
}

// LikePost documentation
// @Summary Like Post
// @Description When user like post
// @Tags like_post
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID to create like post"
// @Security ApiKeyAuth
// @Router /posts/like_post/{post_id} [post]
//func (p *PostLikeController) LikePost(ctx *gin.Context) {
//	postIdStr := ctx.Param("post_id")
//	postId, err := uuid.Parse(postIdStr)
//	if err != nil {
//		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
//		return
//	}
//
//	userIdClaim, err := extensions.GetUserID(ctx)
//	if err != nil {
//		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
//		return
//	}
//	likeUserPostModel := mapper2.MapToLikeUserPostFromPostIdAndUserId(postId, userIdClaim)
//	postDto, resultCode, httpStatusCode, err := services.LikeUserPost().LikePost(ctx, likeUserPostModel, userIdClaim)
//	if err != nil {
//		pkg_response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
//		return
//	}
//
//	pkg_response.SuccessResponse(ctx, pkg_response.ErrCodeSuccess, httpStatusCode, postDto)
//}

// GetUserLikePost documentation
// @Summary Get User like posts
// @Description Retrieve multiple posts filtered by various criteria.
// @Tags like_post
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID to get user like post"
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Security ApiKeyAuth
// @Router /posts/like_post/{post_id} [get]
//func (p *PostLikeController) GetUserLikePost(ctx *gin.Context) {
//	postIdStr := ctx.Param("post_id")
//	postId, err := uuid.Parse(postIdStr)
//	if err != nil {
//		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, "Invalid postId format")
//		return
//	}
//
//	var query query.PostLikeQueryObject
//	if err := ctx.ShouldBindQuery(&query); err != nil {
//		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, "Invalid query parameters")
//		return
//	}
//
//	likeUserPost, resultCode, httpStatusCode, paging, err := services.LikeUserPost().GetUsersOnLikes(ctx, postId, &query)
//	if err != nil {
//		pkg_response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
//		return
//	}
//
//	var userDtos []response.UserDtoShortVer
//	for _, user := range likeUserPost {
//		userDto := mapper.MapUserToUserDtoShortVer(user)
//		userDtos = append(userDtos, userDto)
//	}
//
//	pkg_response.SuccessPagingResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, userDtos, *paging)
//}
