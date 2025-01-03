package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	userQuery "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	userRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
	"net/http"
)

type sUserFriend struct {
	userRepo          userRepo.IUserRepository
	friendRequestRepo userRepo.IFriendRequestRepository
	friendRepo        userRepo.IFriendRepository
	notificationRepo  userRepo.INotificationRepository
}

func NewUserFriendImplement(
	userRepo userRepo.IUserRepository,
	friendRequestRepo userRepo.IFriendRequestRepository,
	friendRepo userRepo.IFriendRepository,
	notificationRepo userRepo.INotificationRepository,
) *sUserFriend {
	return &sUserFriend{
		userRepo:          userRepo,
		friendRequestRepo: friendRequestRepo,
		friendRepo:        friendRepo,
		notificationRepo:  notificationRepo,
	}
}

func (s *sUserFriend) SendAddFriendRequest(
	ctx context.Context,
	command *userCommand.SendAddFriendRequestCommand,
) (result *userCommand.SendAddFriendRequestCommandResult, err error) {
	result = &userCommand.SendAddFriendRequestCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Check exist friend
	friendEntity, err := userEntity.NewFriend(command.UserId, command.FriendId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendCheck, err := s.friendRepo.CheckFriendExist(ctx, friendEntity)
	if err != nil {
		return result, fmt.Errorf("failed to check friend: %w", err)
	}

	// 2. Return if friend has already exist
	if friendCheck {
		result.ResultCode = response.ErrFriendHasAlreadyExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("friend has already exist, you don't need to request more")
	}

	// 3. Find exist friends request
	friendRequestEntityFromUserFound, err := userEntity.NewFriendRequest(command.FriendId, command.UserId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendRequestFromUserFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFromUserFound)
	if err != nil {
		return result, fmt.Errorf("failed to check friend: %w", err)
	}

	if friendRequestFromUserFound {
		result.ResultCode = response.ErrFriendHasAlreadyExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("your friend has already send add friend request, you don't need to request more")
	}

	friendRequestEntityFound, err := userEntity.NewFriendRequest(command.UserId, command.FriendId)

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFound)
	if err != nil {
		return result, fmt.Errorf("failed to check friend request: %w", err)
	}

	// 4. Return if friend request has already exist
	if friendRequestFound {
		result.ResultCode = response.ErrFriendHasAlreadyExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("friend request already exists, you don't need to request more")
	}

	// 5. Find user and friend
	userFound, err := s.userRepo.GetOne(ctx, "id=?", command.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("user record not found: %w", err)
		}
		return result, fmt.Errorf("failed to get user: %w", err)
	}

	friendFound, err := s.userRepo.GetOne(ctx, "id=?", command.FriendId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("friend record not found: %w", err)
		}
		return result, fmt.Errorf("failed to get friend: %w", err)
	}

	// 6. Create friend request
	err = s.friendRequestRepo.CreateOne(ctx, &userEntity.FriendRequest{
		UserId:   command.UserId,
		FriendId: command.FriendId,
	})
	if err != nil {
		return result, fmt.Errorf("failed to add friend request: %w", err)
	}

	// 7. Push notification to user
	notification, err := notificationEntity.NewNotification(
		userFound.FamilyName+" "+userFound.Name,
		userFound.AvatarUrl,
		friendFound.ID,
		consts.FRIEND_REQUEST,
		userFound.ID.String(),
		"",
	)

	_, err = s.notificationRepo.CreateOne(ctx, notification)
	if err != nil {
		return result, fmt.Errorf("failed to add notification: %w", err)
	}

	// 8. Send realtime notification (websocket)
	userSocketResponse := &consts.UserSocketResponse{
		ID:         friendFound.ID,
		FamilyName: friendFound.FamilyName,
		Name:       friendFound.Name,
		AvatarUrl:  friendFound.AvatarUrl,
	}

	notificationSocketResponse := &consts.NotificationSocketResponse{
		From:             userFound.FamilyName + " " + userFound.Name,
		FromUrl:          userFound.AvatarUrl,
		UserId:           friendFound.ID,
		User:             *userSocketResponse,
		NotificationType: consts.FRIEND_REQUEST,
		ContentId:        (userFound.ID).String(),
	}

	err = global.SocketHub.SendNotification(friendFound.ID.String(), notificationSocketResponse)
	if err != nil {
		return result, fmt.Errorf("failed to send notification: %w", err)
	}

	// 9. Response success
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserFriend) GetFriendRequests(
	ctx context.Context,
	query *userQuery.FriendRequestQuery,
) (result *userQuery.FriendRequestQueryResult, err error) {
	result = &userQuery.FriendRequestQueryResult{
		Users:          nil,
		PagingResponse: nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Get list of user request to add friend
	userEntities, paging, err := s.friendRequestRepo.GetFriendRequests(ctx, query)
	if err != nil {
		return result, fmt.Errorf("failed to get friend requests: %w", err)
	}

	// 2. Map userEntity to userDtoShortVer
	var userResults []*common.UserShortVerResult
	for _, user := range userEntities {
		userResults = append(userResults, mapper.NewUserShortVerEntity(user))
	}

	result.Users = userResults
	result.PagingResponse = paging
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserFriend) AcceptFriendRequest(
	ctx context.Context,
	command *userCommand.AcceptFriendRequestCommand,
) (result *userCommand.AcceptFriendRequestCommandResult, err error) {
	result = &userCommand.AcceptFriendRequestCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Find exist friends request
	friendRequestEntityFound, err := userEntity.NewFriendRequest(command.UserId, command.FriendId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFound)
	if err != nil {
		return result, fmt.Errorf("failed to check friend request: %w", err)
	}

	// 2. Return if friend request is not exist
	if !friendRequestFound {
		result.ResultCode = response.ErrFriendNotExist
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("friend request is not exist: %w", err)
	}

	// 3. Find user and friend
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("user record not found: %w", err)
		}
		return result, fmt.Errorf("failed to get user: %w", err)
	}

	friendFound, err := s.userRepo.GetById(ctx, command.FriendId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("friend record not found: %w", err)
		}
		return result, fmt.Errorf("failed to get friend: %w", err)
	}

	// 4. Create friend
	friendEntityForUser, err := userEntity.NewFriend(userFound.ID, friendFound.ID)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendEntityForFriend, err := userEntity.NewFriend(friendFound.ID, userFound.ID)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	err = s.friendRepo.CreateOne(ctx, friendEntityForUser)
	if err != nil {
		return result, fmt.Errorf("failed to add friend: %w", err)
	}

	err = s.friendRepo.CreateOne(ctx, friendEntityForFriend)
	if err != nil {
		return result, fmt.Errorf("failed to add friend: %w", err)
	}

	// 5. Delete friendRequest
	err = s.friendRequestRepo.DeleteOne(ctx, &userEntity.FriendRequest{
		UserId:   command.UserId,
		FriendId: command.FriendId,
	})
	if err != nil {
		return result, fmt.Errorf("failed to delete friend request: %w", err)
	}

	// 6. Plus +1 to friend count for user and friend
	updateUserData := &userEntity.UserUpdate{
		FriendCount: pointer.Ptr(userFound.FriendCount + 1),
	}

	updateFriendData := &userEntity.UserUpdate{
		FriendCount: pointer.Ptr(friendFound.FriendCount + 1),
	}

	err = updateUserData.ValidateUserUpdate()
	if err != nil {
		return result, fmt.Errorf("failed to update user: %w", err)
	}

	err = updateFriendData.ValidateUserUpdate()
	if err != nil {
		return result, fmt.Errorf("failed to update friend: %w", err)
	}

	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, updateUserData)
	if err != nil {
		return result, fmt.Errorf("failed to update user: %w", err)
	}

	_, err = s.userRepo.UpdateOne(ctx, friendFound.ID, updateFriendData)
	if err != nil {
		return result, fmt.Errorf("failed to update friend: %w", err)
	}

	// 7. Push notification to user
	notification := &notificationEntity.Notification{
		From:             friendFound.FamilyName + " " + friendFound.Name,
		FromUrl:          friendFound.AvatarUrl,
		UserId:           userFound.ID,
		NotificationType: consts.ACCEPT_FRIEND_REQUEST,
		ContentId:        (friendFound.ID).String(),
	}

	_, err = s.notificationRepo.CreateOne(ctx, notification)
	if err != nil {
		return result, fmt.Errorf("failed to add notification: %w", err)
	}

	// 8. Send realtime notification (websocket)
	userSocketResponse := &consts.UserSocketResponse{
		ID:         userFound.ID,
		FamilyName: userFound.FamilyName,
		Name:       userFound.Name,
		AvatarUrl:  userFound.AvatarUrl,
	}

	notificationSocketResponse := &consts.NotificationSocketResponse{
		From:             friendFound.FamilyName + " " + friendFound.Name,
		FromUrl:          friendFound.AvatarUrl,
		UserId:           userFound.ID,
		User:             *userSocketResponse,
		NotificationType: consts.FRIEND_REQUEST,
		ContentId:        (friendFound.ID).String(),
	}

	err = global.SocketHub.SendNotification(userFound.ID.String(), notificationSocketResponse)
	if err != nil {
		return result, fmt.Errorf("failed to send notification: %w", err)
	}

	// 9. Response success
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserFriend) RemoveFriendRequest(
	ctx context.Context,
	command *userCommand.RemoveFriendRequestCommand,
) (result *userCommand.RemoveFriendRequestCommandResult, err error) {
	result = &userCommand.RemoveFriendRequestCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Find exist friends request
	friendRequestEntityFound, err := userEntity.NewFriendRequest(command.UserId, command.FriendId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFound)
	if err != nil {
		return result, fmt.Errorf("failed to check friend request: %w", err)
	}

	// 2. Return if friend request is not exist
	if !friendRequestFound {
		result.ResultCode = response.ErrFriendRequestNotExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("friend request is not exist: %w", err)
	}

	// 3. Delete friend request
	friendRequestEntity, err := userEntity.NewFriendRequest(command.UserId, command.FriendId)

	err = s.friendRequestRepo.DeleteOne(ctx, friendRequestEntity)
	if err != nil {
		return result, fmt.Errorf("failed to delete friend request: %w", err)
	}

	// 4. Response success
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusNoContent
	return result, nil
}

func (s *sUserFriend) UnFriend(
	ctx context.Context,
	command *userCommand.UnFriendCommand,
) (result *userCommand.UnFriendCommandResult, err error) {
	result = &userCommand.UnFriendCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Check friend exist
	friendEntity, err := userEntity.NewFriend(command.UserId, command.FriendId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendCheck, err := s.friendRepo.CheckFriendExist(ctx, friendEntity)
	if err != nil {
		return result, fmt.Errorf("failed to check friend: %w", err)
	}

	if !friendCheck {
		result.ResultCode = response.ErrFriendNotExist
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("friend is not exist: %w", err)
	}

	// 2. Remove friend
	err = s.friendRepo.DeleteOne(ctx, friendEntity)
	if err != nil {
		return result, fmt.Errorf("failed to delete friend: %w", err)
	}

	friendEntityForFriend, err := userEntity.NewFriend(command.FriendId, command.UserId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	err = s.friendRepo.DeleteOne(ctx, friendEntityForFriend)
	if err != nil {
		return result, fmt.Errorf("failed to delete friend: %w", err)
	}

	// 3. Minus -1 to friend count of user and friend
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		return result, fmt.Errorf("failed to get user: %w", err)
	}

	friendFound, err := s.userRepo.GetById(ctx, friendEntityForFriend.UserId)
	if err != nil {
		return result, fmt.Errorf("failed to get friend: %w", err)
	}

	updateUserData := &userEntity.UserUpdate{
		FriendCount: pointer.Ptr(userFound.FriendCount - 1),
	}

	updateFriendData := &userEntity.UserUpdate{
		FriendCount: pointer.Ptr(friendFound.FriendCount - 1),
	}

	err = updateUserData.ValidateUserUpdate()
	if err != nil {
		return result, fmt.Errorf("failed to update user: %w", err)
	}

	err = updateFriendData.ValidateUserUpdate()
	if err != nil {
		return result, fmt.Errorf("failed to update friend: %w", err)
	}

	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, updateUserData)
	if err != nil {
		return result, fmt.Errorf("failed to update user: %w", err)
	}

	_, err = s.userRepo.UpdateOne(ctx, friendFound.ID, updateFriendData)
	if err != nil {
		return result, fmt.Errorf("failed to update friend: %w", err)
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserFriend) GetFriends(
	ctx context.Context,
	query *userQuery.FriendQuery,
) (result *userQuery.FriendQueryResult, err error) {
	result = &userQuery.FriendQueryResult{
		Users:          nil,
		PagingResponse: nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Get list of friend
	userEntities, paging, err := s.friendRepo.GetFriends(ctx, query)
	if err != nil {
		return result, fmt.Errorf("failed to get friends: %w", err)
	}

	// 2. Map userModel to userResultShortVer
	var userResults []*common.UserShortVerResult
	for _, user := range userEntities {
		userResults = append(userResults, mapper.NewUserShortVerEntity(user))
	}

	result.Users = userResults
	result.PagingResponse = paging
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
