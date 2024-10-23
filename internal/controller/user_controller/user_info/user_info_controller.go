package user_info

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/user_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"mime/multipart"
	"net/http"
)

type cUserInfo struct{}

func NewUserInfoController() *cUserInfo {
	return &cUserInfo{}
}

// GetInfoByUserId documentation
// @Summary Get user by ID
// @Description Retrieve a user by its unique ID
// @Tags user_info
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

	user, resultCode, httpStatusCode, err := services.UserInfo().GetInfoByUserId(ctx, userId)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	userDto := mapper.MapUserToUserDtoWithoutSetting(user)

	response.SuccessResponse(ctx, resultCode, http.StatusOK, userDto)
}

// GetManyUsers documentation
// @Summary      Get a list of users
// @Description  Retrieve users based on filters such as name, email, phone number, birthday, and created date. Supports pagination and sorting.
// @Tags         user_info
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

	users, resultCode, httpStatusCode, paging, err := services.UserInfo().GetManyUsers(ctx, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	var userDtos []user_dto.UserDtoWithoutSetting
	for _, user := range users {
		userDto := mapper.MapUserToUserDtoWithoutSetting(user)
		userDtos = append(userDtos, *userDto)
	}

	response.SuccessPagingResponse(ctx, resultCode, http.StatusOK, userDtos, *paging)
}

// UpdateUser godoc
// @Summary      Update user information
// @Description  Update various fields of the user profile including name, email, phone number, birthday, and upload avatar and capwall images.
// @Tags         user_info
// @Accept       multipart/form-data
// @Produce      json
// @Param        family_name      formData  string  false  "User's family name"
// @Param        name             formData  string  false  "User's given name"
// @Param        email            formData  string  false  "User's email address"
// @Param        phone_number     formData  string  false  "User's phone number"
// @Param        birthday         formData  string  false  "User's birthday"
// @Param        avatar_url       formData  file    false  "Upload user avatar image"
// @Param        capwall_url      formData  file    false  "Upload user capwall image"
// @Param        privacy          formData  string  false  "User privacy level"
// @Param        biography        formData  string  false  "User biography"
// @Param        language_setting formData  string  false  "Setting language "vi" or "en""
// @Success      200              {object}  response.ResponseData
// @Failure      500              {object}  response.ErrResponse
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

	var openFileAvatar multipart.File
	var openFileCapwall multipart.File

	if updateInput.AvatarUrl.Size != 0 {
		openFileAvatar, err = updateInput.AvatarUrl.Open()
		if err != nil {
			response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if updateInput.CapwallUrl.Size != 0 {
		openFileCapwall, err = updateInput.CapwallUrl.Open()
		if err != nil {
			response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
	}

	var languageSetting consts.Language
	if updateInput.LanguageSetting != nil {
		languageSetting = *updateInput.LanguageSetting
	}

	user, resultCode, httpStatusCode, err := services.UserInfo().UpdateUser(ctx, userIdClaim, updateData, openFileAvatar, openFileCapwall, languageSetting)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
	}

	userDto := mapper.MapUserToUserDto(user)

	response.SuccessResponse(ctx, resultCode, http.StatusOK, userDto)
}
