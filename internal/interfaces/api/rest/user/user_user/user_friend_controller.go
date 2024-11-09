package user_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/user/user_user/dto/mapper"
	query2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/user/user_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cUserFriend struct{}

func NewUserFriendController() *cUserFriend {
	return &cUserFriend{}
}

// SendAddFriendRequest godoc
// @Summary Send add friend request
// @Description Send add friend request to another people
// @Tags user_friend
// @Param friend_id path string true "User id you want to send add request"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/friends/friend_request/{friend_id}/ [post]
func (c *cUserFriend) SendAddFriendRequest(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Check user send add friend request for himself
	if userIdClaim == friendId {
		response.ErrorResponse(ctx, response.ErrMakeFriendWithYourSelf, http.StatusBadRequest, "You can not make friend with yourself")
		return
	}

	// 3. Map to friendRequestModel
	friendRequestModel := mapper.MapToFriendRequestFromUserIdAndFriendId(userIdClaim, friendId)

	// 4. Call service
	resultCode, httpStatusCode, err := services.UserFriend().SendAddFriendRequest(ctx, friendRequestModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, nil)
}

// UndoFriendRequest godoc
// @Summary Undo add friend request
// @Description Undo add friend request
// @Tags user_friend
// @Param friend_id path string true "User id you want to undo add request"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/friends/friend_request/{friend_id}/ [delete]
func (c *cUserFriend) UndoFriendRequest(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Map to friendRequestModel
	friendRequestModel := mapper.MapToFriendRequestFromUserIdAndFriendId(userIdClaim, friendId)

	// 4. Call service
	resultCode, httpStatusCode, err := services.UserFriend().RemoveFriendRequest(ctx, friendRequestModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, nil)
}

// GetFriendRequests godoc
// @Summary Get a list of friend request
// @Description Get a list of friend request
// @Tags user_friend
// @Param limit query int false "limit on page"
// @Param page query int false "current page"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/friends/friend_request [get]
func (c *cUserFriend) GetFriendRequests(ctx *gin.Context) {
	// 1. Validate and get query object from query
	var query query2.FriendRequestQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call services
	userDtos, paging, resultCode, httpStatusCode, err := services.UserFriend().GetFriendRequests(ctx, userIdClaim, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessPagingResponse(ctx, resultCode, http.StatusOK, userDtos, *paging)
}

// AcceptFriendRequest godoc
// @Summary Accept friend request
// @Description Accept friend request
// @Tags user_friend
// @Param friend_id path string true "User id you want to accept friend request"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/friends/friend_response/{friend_id}/ [post]
func (c *cUserFriend) AcceptFriendRequest(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Map to friendRequestModel
	friendRequestModel := mapper.MapToFriendRequestFromUserIdAndFriendId(friendId, userIdClaim)

	// 4. Call service
	resultCode, httpStatusCode, err := services.UserFriend().AcceptFriendRequest(ctx, friendRequestModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, nil)
}

// RejectFriendRequest godoc
// @Summary Reject friend request
// @Description Delete friend request
// @Tags user_friend
// @Param friend_id path string true "User id you want to reject friend request"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/friends/friend_response/{friend_id}/ [delete]
func (c *cUserFriend) RejectFriendRequest(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Map to friendRequestModel
	friendRequestModel := mapper.MapToFriendRequestFromUserIdAndFriendId(friendId, userIdClaim)

	// 4. Call service
	resultCode, httpStatusCode, err := services.UserFriend().RemoveFriendRequest(ctx, friendRequestModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, nil)
}

// UnFriend godoc
// @Summary unfriend
// @Description unfriend
// @Tags user_friend
// @Param friend_id path string true "User id you want to unfriend"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/friends/{friend_id}/ [delete]
func (c *cUserFriend) UnFriend(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Map to friendRequestModel
	friendModel := mapper.MapToFriendFromUserIdAndFriendId(friendId, userIdClaim)

	// 4. Call service
	resultCode, httpStatusCode, err := services.UserFriend().UnFriend(ctx, friendModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, nil)
}

// GetFriends godoc
// @Summary Get a list of friend
// @Description Get a list of friend
// @Tags user_friend
// @Param limit query int false "limit on page"
// @Param page query int false "current page"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/friends/ [get]
func (c *cUserFriend) GetFriends(ctx *gin.Context) {
	// 1. Validate and get query object from query
	var query query2.FriendQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call services
	userDtos, paging, resultCode, httpStatusCode, err := services.UserFriend().GetFriends(ctx, userIdClaim, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessPagingResponse(ctx, resultCode, http.StatusOK, userDtos, *paging)
}
