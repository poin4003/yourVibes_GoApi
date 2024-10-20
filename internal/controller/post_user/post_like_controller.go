package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/redis/go-redis/v9"
	"net/http"
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
// @Tags post
// @Accept json
// @Produce json
// @Param postId path string true "Post ID to create like post"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/like_post/{postId} [post]
func (p *PostLikeController) LikePost(ctx *gin.Context) {
	postIdStr := ctx.Param("postId")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userUUID, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}
	likeUserPostModel := mapper.MapToLikeUserPostFromPostIdAndUserId(postId, userUUID)
	resultCode, err := services.LikeUserPost().LikePost(ctx, likeUserPostModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, likeUserPostModel)
}

// DeleteLikePost documentation
// @Summary Delete Post
// @Description When user delete like post
// @Tags post
// @Accept json
// @Produce json
// @Param postId path string true "Post ID to delete like post"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/delete_post/{postId} [delete]
func (p *PostLikeController) DeleteLikePost(ctx *gin.Context) {
	postIdStr := ctx.Param("postId")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	deleteLikeUserPostModel := mapper.MapToLikeUserPostFromPostIdAndUserId(postId, userIdClaim)
	resultCode, err := services.LikeUserPost().DeleteLikePost(ctx, deleteLikeUserPostModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, deleteLikeUserPostModel)
}

// GetUserLikePost documentation
// @Summary Get User like posts
// @Description Retrieve multiple posts filtered by various criteria.
// @Tags post
// @Accept json
// @Produce json
// @Param postId path string true "Post ID to get user like post"
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /posts/likeUser/{postId} [get]
func (p *PostLikeController) GetUserLikePost(ctx *gin.Context) {
	// Lấy postId từ URL và kiểm tra tính hợp lệ
	postIdStr := ctx.Param("postId")
	postId, err := uuid.Parse(postIdStr) // Kiểm tra postId có đúng định dạng UUID không
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, "Invalid postId format")
		return
	}

	// Ràng buộc query parameters với PostLikeQueryObject
	var query query_object.PostLikeQueryObject
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Đặt giá trị mặc định cho pagination nếu không được cung cấp
	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	query.PostID = postId.String() // Gán postId vào query object

	// Gọi service để lấy danh sách user đã like
	likeUserPost, resultCode, err := services.LikeUserPost().GetUsersOnLikes(ctx, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	// Tạo response với dữ liệu phân trang
	total := int64(len(likeUserPost))
	paging := response.PagingResponse{
		Total: total,
		Page:  query.Page,
		Limit: query.Limit,
	}

	// Trả về kết quả thành công
	response.SuccessPagingResponse(ctx, response.ErrCodeSuccess, http.StatusOK, likeUserPost, paging)
}
