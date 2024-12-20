package comment_admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_admin/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_admin/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cAdminCommentReport struct{}

func NewAdminCommentReportController() *cAdminCommentReport {
	return &cAdminCommentReport{}
}

// GetCommentReport documentation
// @Summary Get comment report detail
// @Description Retrieve a comment report
// @Tags admin_comment_report
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param reported_comment_id path string true "Reported comment id"
// @Security ApiKeyAuth
// @Router /comments/report/{user_id}/{reported_comment_id} [get]
func (c *cAdminCommentReport) GetCommentReport(ctx *gin.Context) {
	// 1. Get userId from param path
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get reportedCommentId from param path
	reportedCommentIdStr := ctx.Param("reported_comment_id")
	reportedCommentId, err := uuid.Parse(reportedCommentIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Call service to handle get comment report
	getOneCommentReportQuery, err := query.ToGetOneCommentReportQuery(userId, reportedCommentId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.CommentReport().GetDetailCommentReport(ctx, getOneCommentReportQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	commentReportDto := response.ToCommentReportDto(result.CommentReport)

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, commentReportDto)
}

// GetManyCommentReports godoc
// @Summary      Get a list of comment report
// @Description  Retrieve comment report base on filters
// @Tags         admin_comment_report
// @Accept       json
// @Produce      json
// @Param        reason        query     string  false  "reason to filter report"
// @Param        status        query     bool    false  "Filter by status"
// @Param        created_at    query     string  false  "Filter by creation day"
// @Param        sort_by       query     string  false  "Sort by field"
// @Param        isDescending  query     bool    false  "Sort in descending order"
// @Param        limit         query     int     false  "Number of results per page"
// @Param        page          query     int     false  "Page number"
// @Security ApiKeyAuth
// @Router       /comments/report [get]
func (c *cAdminCommentReport) GetManyCommentReports(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to userQueryObject
	commentReportQueryObject, ok := queryInput.(*query.CommentReportQueryObject)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3 Call service to handle get many
	getManyCommentReportQuery, err := commentReportQueryObject.ToGetManyCommentQuery()
	result, err := services.CommentReport().GetManyCommentReport(ctx, getManyCommentReportQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	var commentReportDtos []*response.CommentReportShortVerDto
	for _, commentReportResult := range result.CommentReports {
		commentReportDtos = append(commentReportDtos, response.ToCommentReportShortVerDto(commentReportResult))
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, result.HttpStatusCode, commentReportDtos, *result.PagingResponse)
}

// HandleCommentReport godoc
// @Summary handle comment report
// @Description When admin need to handle report
// @Tags admin_comment_report
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param reported_comment_id path string true "Reported comment id"
// @Security ApiKeyAuth
// @Router /comments/report/{user_id}/{reported_comment_id} [patch]
func (c *cAdminCommentReport) HandleCommentReport(ctx *gin.Context) {
	// 1. Get userId from param path
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get reportedCommentId from param path
	reportedCommentIdStr := ctx.Param("reported_comment_id")
	reportedCommentId, err := uuid.Parse(reportedCommentIdStr)
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

	// 4. Call service to handle comment report
	handleCommentReportCommand, err := request.ToHandleCommentReportCommand(adminIdClaim, userId, reportedCommentId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.CommentReport().HandleCommentReport(ctx, handleCommentReportCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. response
	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}
