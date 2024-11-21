package user_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
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
// @Security ApiKeyAuth
// @Router /users/{userId} [get]
func (c *cUserInfo) GetInfoByUserId(ctx *gin.Context) {
	var userRequest query.UserQueryObject

	// 1. Get userId from param path
	userIdStr := ctx.Param("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get userId from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call services
	getOneUserQuery, err := userRequest.ToGetOneUserQuery(userId, userIdClaim)

	result, err := services.UserInfo().GetInfoByUserId(ctx, getOneUserQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	userDto := response.ToUserWithoutSettingDto(result.User)

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, userDto)
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
// @Security ApiKeyAuth
// @Router       /users/ [get]
func (c *cUserInfo) GetManyUsers(ctx *gin.Context) {
	// 1. Get query
	var query query.UserQueryObject
	if err := ctx.ShouldBindQuery(&query); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Call service to handle get many
	getManyUserQuery, err := query.ToGetManyUserQuery()
	result, err := services.UserInfo().GetManyUsers(ctx, getManyUserQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 3. Map to dto
	var userDtos []*response.UserShortVerDto
	for _, userResult := range result.Users {
		userDtos = append(userDtos, response.ToUserShortVerDto(userResult))
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, http.StatusOK, userDtos, *result.PagingResponse)
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
// @Security ApiKeyAuth
// @Router       /users/ [patch]
func (*cUserInfo) UpdateUser(ctx *gin.Context) {
	// 1. Get body from form
	var updateInput request.UpdateUserRequest
	if err := ctx.ShouldBind(&updateInput); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Get avatar and capwall from body
	var openFileAvatar multipart.File
	var openFileCapwall multipart.File

	if updateInput.Avatar.Size != 0 {
		openFileAvatar, err = updateInput.Avatar.Open()
		if err != nil {
			pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if updateInput.Capwall.Size != 0 {
		openFileCapwall, err = updateInput.Capwall.Open()
		if err != nil {
			pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// 4. Call service to handle update user
	updateUserCommand, err := updateInput.ToUpdateUserCommand(userIdClaim, openFileAvatar, openFileCapwall)
	result, err := services.UserInfo().UpdateUser(ctx, updateUserCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
	}

	// 5. Map to dto
	userDto := response.ToUserWithSettingDto(result.User)

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, userDto)
}
