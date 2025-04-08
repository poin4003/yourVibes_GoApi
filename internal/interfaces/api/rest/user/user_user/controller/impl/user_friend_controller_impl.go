package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
)

type cUserFriend struct {
	userFriendService services.IUserFriend
}

func NewUserFriendController(
	userFriendService services.IUserFriend,
) *cUserFriend {
	return &cUserFriend{
		userFriendService: userFriendService,
	}
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
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 3. Check user send add friend request for himself
	if userIdClaim == friendId {
		ctx.Error(pkgResponse.NewCustomError(pkgResponse.ErrMakeFriendWithYourSelf))
		return
	}

	// 5. Call service
	sendFriendRequestCommand := &command.SendAddFriendRequestCommand{
		UserId:   userIdClaim,
		FriendId: friendId,
	}

	err = c.userFriendService.SendAddFriendRequest(ctx, sendFriendRequestCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
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
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 3. Call service
	removeFriendRequestCommand := &command.RemoveFriendRequestCommand{
		UserId:   userIdClaim,
		FriendId: friendId,
	}

	err = c.userFriendService.RemoveFriendRequest(ctx, removeFriendRequestCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
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
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to userQueryObject
	friendRequestQueryObject, ok := queryInput.(*query.FriendRequestQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call services
	friendRequestQuery, _ := friendRequestQueryObject.ToFriendRequestQuery(userIdClaim)
	result, err := c.userFriendService.GetFriendRequests(ctx, friendRequestQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	var userDtos []*response.UserShortVerDto
	for _, userResult := range result.Users {
		userDtos = append(userDtos, response.ToUserShortVerDto(userResult))
	}

	pkgResponse.OKWithPaging(ctx, userDtos, *result.PagingResponse)
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
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 3. Call service
	friendRequestCommand := &command.AcceptFriendRequestCommand{
		UserId:   friendId,
		FriendId: userIdClaim,
	}

	err = c.userFriendService.AcceptFriendRequest(ctx, friendRequestCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
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
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 3. Call service
	friendRequestCommand := &command.RemoveFriendRequestCommand{
		UserId:   friendId,
		FriendId: userIdClaim,
	}
	err = c.userFriendService.RemoveFriendRequest(ctx, friendRequestCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
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
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get user id claim from jwt
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 3. Call service
	unFriendCommand := &command.UnFriendCommand{
		UserId:   friendId,
		FriendId: userIdClaim,
	}
	err = c.userFriendService.UnFriend(ctx, unFriendCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
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
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to userQueryObject
	friendQueryObject, ok := queryInput.(*query.FriendQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get user id from param
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 4. Call services
	friendQuery, _ := friendQueryObject.ToFriendQuery(userId)

	result, err := c.userFriendService.GetFriends(ctx, friendQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	var userDtos []*response.UserShortVerDto
	for _, userResult := range result.Users {
		userDtos = append(userDtos, response.ToUserShortVerDto(userResult))
	}

	pkgResponse.OKWithPaging(ctx, userDtos, *result.PagingResponse)
}

// GetFriendSuggestion godoc
// @Summary Get a list of friend suggestion
// @Description Get a list of friend
// @Tags user_friend
// @Param limit query int false "limit on page"
// @Param page query int false "current page"
// @Security ApiKeyAuth
// @Router /users/friends/suggestion/ [get]
func (c *cUserFriend) GetFriendSuggestion(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to userQueryObject
	friendQueryObject, ok := queryInput.(*query.FriendQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get user id from param
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call services
	friendQuery, _ := friendQueryObject.ToFriendQuery(userIdClaim)

	result, err := c.userFriendService.GetFriendSuggestion(ctx, friendQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	var userDtos []*response.UserShortVerWithFriendSuggestionDto
	for _, userResult := range result.Users {
		userDtos = append(userDtos, response.ToUserShortWithSendFriendRequestVerDto(userResult))
	}

	pkgResponse.OKWithPaging(ctx, userDtos, *result.PagingResponse)
}
