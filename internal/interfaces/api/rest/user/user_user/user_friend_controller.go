package user_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
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
// @Security ApiKeyAuth
// @Router /users/friends/friend_request/{friend_id}/ [post]
func (c *cUserFriend) SendAddFriendRequest(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Check user send add friend request for himself
	if userIdClaim == friendId {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrMakeFriendWithYourSelf, http.StatusBadRequest, "You can not make friend with yourself")
		return
	}

	// 5. Call service
	sendFriendRequestCommand := &command.SendAddFriendRequestCommand{
		UserId:   userIdClaim,
		FriendId: friendId,
	}

	result, err := services.UserFriend().SendAddFriendRequest(ctx, sendFriendRequestCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, nil)
}

// UndoFriendRequest godoc
// @Summary Undo add friend request
// @Description Undo add friend request
// @Tags user_friend
// @Param friend_id path string true "User id you want to undo add request"
// @Security ApiKeyAuth
// @Router /users/friends/friend_request/{friend_id}/ [delete]
func (c *cUserFriend) UndoFriendRequest(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call service
	removeFriendRequestCommand := &command.RemoveFriendRequestCommand{
		UserId:   userIdClaim,
		FriendId: friendId,
	}

	result, err := services.UserFriend().RemoveFriendRequest(ctx, removeFriendRequestCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, nil)
}

// GetFriendRequests godoc
// @Summary Get a list of friend request
// @Description Get a list of friend request
// @Tags user_friend
// @Param limit query int false "limit on page"
// @Param page query int false "current page"
// @Security ApiKeyAuth
// @Router /users/friends/friend_request [get]
func (c *cUserFriend) GetFriendRequests(ctx *gin.Context) {
	// 1. Validate and get query object from query
	var query query.FriendRequestQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call services
	friendRequestQuery, err := query.ToFriendRequestQuery(userIdClaim)
	result, err := services.UserFriend().GetFriendRequests(ctx, friendRequestQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	var userDtos []*response.UserShortVerDto
	for _, userResult := range result.Users {
		userDtos = append(userDtos, response.ToUserShortVerDto(userResult))
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, http.StatusOK, userDtos, *result.PagingResponse)
}

// AcceptFriendRequest godoc
// @Summary Accept friend request
// @Description Accept friend request
// @Tags user_friend
// @Param friend_id path string true "User id you want to accept friend request"
// @Security ApiKeyAuth
// @Router /users/friends/friend_response/{friend_id}/ [post]
func (c *cUserFriend) AcceptFriendRequest(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call service
	friendRequestCommand := &command.AcceptFriendRequestCommand{
		UserId:   friendId,
		FriendId: userIdClaim,
	}

	result, err := services.UserFriend().AcceptFriendRequest(ctx, friendRequestCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, nil)
}

// RejectFriendRequest godoc
// @Summary Reject friend request
// @Description Delete friend request
// @Tags user_friend
// @Param friend_id path string true "User id you want to reject friend request"
// @Security ApiKeyAuth
// @Router /users/friends/friend_response/{friend_id}/ [delete]
func (c *cUserFriend) RejectFriendRequest(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call service
	friendRequestCommand := &command.RemoveFriendRequestCommand{
		UserId:   friendId,
		FriendId: userIdClaim,
	}
	result, err := services.UserFriend().RemoveFriendRequest(ctx, friendRequestCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, nil)
}

// UnFriend godoc
// @Summary unfriend
// @Description unfriend
// @Tags user_friend
// @Param friend_id path string true "User id you want to unfriend"
// @Security ApiKeyAuth
// @Router /users/friends/{friend_id}/ [delete]
func (c *cUserFriend) UnFriend(ctx *gin.Context) {
	// 1. Get friend id from param path
	friendIdStr := ctx.Param("friend_id")
	friendId, err := uuid.Parse(friendIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call service
	unFriendCommand := &command.UnFriendCommand{
		UserId:   friendId,
		FriendId: userIdClaim,
	}
	result, err := services.UserFriend().UnFriend(ctx, unFriendCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, nil)
}

// GetFriends godoc
// @Summary Get a list of friend
// @Description Get a list of friend
// @Tags user_friend
// @Param user_id path string true "User id you want to get a friend list"
// @Param limit query int false "limit on page"
// @Param page query int false "current page"
// @Security ApiKeyAuth
// @Router /users/friends/{user_id} [get]
func (c *cUserFriend) GetFriends(ctx *gin.Context) {
	// 1. Validate and get query object from query
	var query query.FriendQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id from param
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Call services
	friendQuery, err := query.ToFriendQuery(userId)

	result, err := services.UserFriend().GetFriends(ctx, friendQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	var userDtos []*response.UserShortVerDto
	for _, userResult := range result.Users {
		userDtos = append(userDtos, response.ToUserShortVerDto(userResult))
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, http.StatusOK, userDtos, *result.PagingResponse)
}
