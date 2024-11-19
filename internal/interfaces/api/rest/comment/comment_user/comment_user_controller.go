package comment_user

type cCommentUser struct {
}

func NewCommentUserController() *cCommentUser {
	return &cCommentUser{}
}

// CreateComment documentation
// @Summary Comment create comment
// @Description When user create comment or rep comment
// @Tags comment_user
// @Accept json
// @Produce json
// @Param input body request.CreateCommentInput true "input"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /comments/ [post]
//func (p *cCommentUser) CreateComment(ctx *gin.Context) {
//	var commentInput request.CreateCommentInput
//
//	if err := ctx.ShouldBindJSON(&commentInput); err != nil {
//		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
//		return
//	}
//
//	userUUID, err := extensions.GetUserID(ctx)
//	if err != nil {
//		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
//		return
//	}
//
//	commentModel := mapper.MapToCommentFromCreateDto(&commentInput, userUUID)
//	comment, resultCode, httpStatusCode, err := services.CommentUser().CreateComment(ctx, commentModel)
//	if err != nil {
//		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
//		return
//	}
//
//	commentDto := mapper.MapCommentToNewCommentDto(comment)
//
//	response.SuccessResponse(ctx, resultCode, http.StatusOK, commentDto)
//}

// GetManyComment documentation
// @Summary Get many comment
// @Description Retrieve multiple comment filtered by various criteria.
// @Tags comment_user
// @Accept json
// @Produce json
// @Param post_id query string true "Post ID to filter comment, get the first layer"
// @Param parent_id query string false "Filter by parent id, get the next layer"
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /comments/ [get]
//func (p *cCommentUser) GetComment(ctx *gin.Context) {
//	var query query.CommentQueryObject
//
//	if err := ctx.ShouldBindQuery(&query); err != nil {
//		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
//		return
//	}
//
//	userUUID, err := extensions.GetUserID(ctx)
//	if err != nil {
//		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
//		return
//	}
//
//	commentDtos, resultCode, httpStatusCode, paging, err := services.CommentUser().GetManyComments(ctx, &query, userUUID)
//	if err != nil {
//		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
//		return
//	}
//
//	response.SuccessPagingResponse(ctx, resultCode, http.StatusOK, commentDtos, *paging)
//}
//
// DeleteComment documentation
// @Summary delete comment by ID
// @Description when user want to delete comment
// @Tags comment_user
// @Accept json
// @Produce json
// @Param comment_id path string true "comment ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /comments/{comment_id} [delete]
//func (p *cCommentUser) DeleteComment(ctx *gin.Context) {
//	commentIdStr := ctx.Param("comment_id")
//	commentId, err := uuid.Parse(commentIdStr)
//	if err != nil {
//		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
//		return
//	}
//
//	resultCode, httpStatusCode, err := services.CommentUser().DeleteComment(ctx, commentId)
//	if err != nil {
//		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
//		return
//	}
//
//	response.SuccessResponse(ctx, resultCode, http.StatusOK, nil)
//}

// UpdateComment documentation
// @Summary update comment
// @Description When user need to update information of comment
// @Tags comment_user
// @Accept multipart/form-data
// @Produce json
// @Param comment_id path string true "commentId"
// @Param input body request.UpdateCommentInput true "input"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /comments/{comment_id} [patch]
//func (p *cCommentUser) UpdateComment(ctx *gin.Context) {
//	var updateInput request.UpdateCommentInput
//
//	if err := ctx.ShouldBindJSON(&updateInput); err != nil {
//		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
//		return
//	}
//
//	commentIdStr := ctx.Param("comment_id")
//	commentId, err := uuid.Parse(commentIdStr)
//	if err != nil {
//		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
//		return
//	}
//
//	commentModel := mapper.MapToCommentFromUpdateDto(&updateInput)
//
//	comment, resultCode, httpStatusCode, err := services.CommentUser().UpdateComment(ctx, commentId, commentModel)
//	if err != nil {
//		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
//		return
//	}
//
//	commentDto := mapper.MapCommentToUpdatedCommentDto(comment)
//
//	response.SuccessResponse(ctx, resultCode, http.StatusOK, commentDto)
//}
