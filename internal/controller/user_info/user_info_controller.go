package user_info

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/user_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cUserInfo struct{}

func NewUserInfoController() *cUserInfo {
	return &cUserInfo{}
}

// GetInfoByUserId documentation
// @Summary Get user by ID
// @Description Retrieve a user by its unique ID
// @Tags user
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/{userId} [get]
func (c *cUserInfo) GetInfoByUserId(ctx *gin.Context) {
	userIdStr := ctx.Param("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	var user *model.User
	var resultCode int

	user, resultCode, err = services.UserInfo().GetInfoByUserId(ctx, userId)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	userDto := mapper.MapUserToUserDto(user)

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, userDto)
}

// GetManyUsers documentation
// @Summary      Get a list of users
// @Description  Retrieve users based on filters such as name, email, phone number, birthday, and created date. Supports pagination and sorting.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        name          query     string  false  "name to filter users"
// @Param        email         query     string  false  "Filter by email"
// @Param        phone_number  query     string  false  "Filter by phone number"
// @Param        birthday      query     string  false  "Filter by birthday"
// @Param        created_at    query     string  false  "Filter by creation day"
// @Param        sort_by       query     string  false  "Sort by field"
// @Param        isDescending  query     bool    false  "Sort in descending order"
// @Param        limit         query     int     false  "Number of results per page"
// @Param        page          query     int     false  "Page number"
// @Success      200           {object}  response.ResponseData
// @Failure      500           {object}  response.ErrResponse
// @Security ApiKeyAuth
// @Router       /users/ [get]
func (c *cUserInfo) GetManyUsers(ctx *gin.Context) {
	var query query_object.UserQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	users, resultCode, err := services.UserInfo().GetManyUsers(ctx, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	total := int64(len(users))

	paging := response.PagingResponse{
		Limit: query.Limit,
		Page:  query.Page,
		Total: total,
	}

	var userDtos []user_dto.UserDto
	for _, user := range users {
		userDto := mapper.MapUserToUserDto(user)
		userDtos = append(userDtos, *userDto)
	}

	response.SuccessPagingResponse(ctx, response.ErrCodeSuccess, http.StatusOK, userDtos, paging)
}

// UpdateUser godoc
// @Summary      Update user information
// @Description  Update various fields of the user profile including name, email, phone number, birthday, and upload avatar and capwall images.
// @Tags         user
// @Accept       multipart/form-data
// @Produce      json
// @Param        family_name   formData  string  false  "User's family name"
// @Param        name          formData  string  false  "User's given name"
// @Param        email         formData  string  false  "User's email address"
// @Param        phone_number  formData  string  false  "User's phone number"
// @Param        birthday      formData  string  false  "User's birthday"
// @Param        avatar_url    formData  file    false  "Upload user avatar image"
// @Param        capwall_url   formData  file    false  "Upload user capwall image"
// @Param        privacy       formData  string  true   "User privacy level"
// @Param        biography     formData  string  false  "User biography"
// @Success      200           {object}  response.ResponseData
// @Failure      500           {object}  response.ErrResponse
// @Security ApiKeyAuth
// @Router       /users/ [patch]
func (*cUserInfo) UpdateUser(ctx *gin.Context) {
	var updateInput user_dto.UpdateUserInput

	if err := ctx.ShouldBind(&updateInput); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	updateData := mapper.MapToUserFromUpdateDto(&updateInput)

	openFileAvatar, err := updateInput.AvatarUrl.Open()
	if err != nil {
		response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	openFileCapwall, err := updateInput.CapwallUrl.Open()
	if err != nil {
		response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	user, resultCode, err := services.UserInfo().UpdateUser(ctx, userIdClaim, updateData, openFileCapwall, openFileAvatar)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
	}

	userDto := mapper.MapUserToUserDto(user)

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, userDto)
}
