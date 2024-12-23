package post_admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_admin/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_admin/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cAdminPostReport struct{}

func NewAdminPostReportController() *cAdminPostReport {
	return &cAdminPostReport{}
}

// GetPostReport documentation
// @Summary Get post report detail
// @Description Retrieve a post report
// @Tags admin_post_report
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param reported_post_id path string true "Reported post id"
// @Security ApiKeyAuth
// @Router /posts/report/{user_id}/{reported_post_id} [get]
func (c *cAdminPostReport) GetPostReport(ctx *gin.Context) {
	// 1. Get userId from param path
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get reportedUserId from param path
	reportedPostIdStr := ctx.Param("reported_post_id")
	reportedPostId, err := uuid.Parse(reportedPostIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Call service to handle get post report
	getOnePostReportQuery, err := query.ToGetOnePostReportQuery(userId, reportedPostId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.PostReport().GetDetailPostReport(ctx, getOnePostReportQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	postReportDto := response.ToPostReportDto(result.PostReport)

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, postReportDto)
}

// GetManyPostReports godoc
// @Summary      Get a list of post report
// @Description  Retrieve post report base on filters
// @Tags         admin_post_report
// @Accept       json
// @Produce      json
// @Param        reason        query     string  false  "reason to filter report"
// @Param        status        query     bool    false  "Filter by status"
// @Param        created_at    query     string  false  "Filter by creation day"
// @Param        user_email    query     string  false  "Filter by user email"
// @Param        admin_email   query     string  false  "Filter by admin email"
// @Param        from_date     query     string  false  "Filter by from date"
// @Param        to_date       query     string  false  "Filter by to date"
// @Param        sort_by       query     string  false  "Sort by field"
// @Param        isDescending  query     bool    false  "Sort in descending order"
// @Param        limit         query     int     false  "Number of results per page"
// @Param        page          query     int     false  "Page number"
// @Security ApiKeyAuth
// @Router       /posts/report [get]
func (c *cAdminPostReport) GetManyPostReports(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to postReportQueryObject
	postReportQueryObject, ok := queryInput.(*query.PostReportQueryObject)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3 Call service to handle get many
	getManyPostReportQuery, err := postReportQueryObject.ToGetManyPostQuery()
	result, err := services.PostReport().GetManyPostReport(ctx, getManyPostReportQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	var postReportDtos []*response.PostReportShortVerDto
	for _, postReportResult := range result.PostReports {
		postReportDtos = append(postReportDtos, response.ToPostReportShortVerDto(postReportResult))
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, result.HttpStatusCode, postReportDtos, *result.PagingResponse)
}

// HandlePostReport godoc
// @Summary handle post report
// @Description When admin need to handle report
// @Tags admin_post_report
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param reported_post_id path string true "Reported post id"
// @Security ApiKeyAuth
// @Router /posts/report/{user_id}/{reported_post_id} [patch]
func (c *cAdminPostReport) HandlePostReport(ctx *gin.Context) {
	// 1. Get userId from param path
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get reportedPostId from param path
	reportedPostIdStr := ctx.Param("reported_post_id")
	reportedPostId, err := uuid.Parse(reportedPostIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Get admin id from token
	adminIdClaim, err := extensions.GetAdminID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Call service to handle post report
	handlePostReportCommand, err := request.ToHandlePostReportCommand(adminIdClaim, userId, reportedPostId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.PostReport().HandlePostReport(ctx, handlePostReportCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. response
	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}

// DeletePostReport godoc
// @Summary delete post report
// @Description When admin need to delete report
// @Tags admin_post_report
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param reported_post_id path string true "Reported post id"
// @Security ApiKeyAuth
// @Router /posts/report/{user_id}/{reported_post_id} [delete]
func (c *cAdminPostReport) DeletePostReport(ctx *gin.Context) {
	// 1. Get userId from param path
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get reportedPostId from param path
	reportedPostIdStr := ctx.Param("reported_post_id")
	reportedPostId, err := uuid.Parse(reportedPostIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Call service to delete post report
	deletePostReportCommand, err := request.ToDeletePostReportCommand(userId, reportedPostId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.PostReport().DeletePostReport(ctx, deletePostReportCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. response
	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}

// ActivatePost godoc
// @Summary activate post account
// @Description When admin need to activate post
// @Tags admin_post_report
// @Accept json
// @Produce json
// @Param post_id path string true "post ID"
// @Security ApiKeyAuth
// @Router /posts/report/activate/{post_id} [patch]
func (c *cAdminPostReport) ActivatePost(ctx *gin.Context) {
	// 1. Get postId from param path
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Call service to activate post
	activatePostCommand := &command.ActivatePostCommand{
		PostId: postId,
	}

	result, err := services.PostReport().ActivatePost(ctx, activatePostCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 3. response
	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}
