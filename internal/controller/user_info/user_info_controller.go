package user_info

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cUserInfo struct {}

var UserInfo = new(cUserInfo)

func (c *cUserInfo) GetInfoByUserId(ctx *gin.Context) {
	userId := ctx.Query("id")

	if userId == "" {
		response.ErrorResponse(ctx, response.NoUserID, http.StatusBadRequest, "user id not found")
		return
	}

	user, err := services.UserInfo().GetInfoByUserId(ctx, userId)
	if err != nil {
		response.ErrorResponse(ctx, response.UserNotFound, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, user)
}

func (c *cUserInfo) GetUsersByName(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	limitStr := ctx.Query("limit")
	pageStr := ctx.Query("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = consts.DEFAULT_LIMIT
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = consts.DEFAULT_PAGE
	}

	if keyword == "" {
		response.ErrorResponse(ctx, response.NoKeywordInFindUsers, http.StatusBadRequest, "keyword is empty")
		return
	}

	users, total, err := services.UserInfo().GetUsersByName(ctx, keyword, limit, page)
	if err != nil {
		response.ErrorResponse(ctx, response.FoundUsersFailed, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessPagingResponse(
		ctx, response.ErrCodeSuccess, http.StatusOK, 
		users, response.PagingResponse{Limit: limit, Page: page, Total: total})
}