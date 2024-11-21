package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notification_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	user_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
)

type sUserFriend struct {
	userRepo          user_repo.IUserRepository
	friendRequestRepo user_repo.IFriendRequestRepository
	friendRepo        user_repo.IFriendRepository
	notificationRepo  user_repo.INotificationRepository
}

func NewUserFriendImplement(
	userRepo user_repo.IUserRepository,
	friendRequestRepo user_repo.IFriendRequestRepository,
	friendRepo user_repo.IFriendRepository,
	notificationRepo user_repo.INotificationRepository,
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
	command *user_command.SendAddFriendRequestCommand,
) (result *user_command.SendAddFriendRequestCommandResult, err error) {
	result = &user_command.SendAddFriendRequestCommandResult{}
	// 1. Check exist friend
	friendEntity, err := user_entity.NewFriend(command.UserId, command.FriendId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendCheck, err := s.friendRepo.CheckFriendExist(ctx, friendEntity)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to check friend: %w", err)
	}

	// 2. Return if friend has already exist
	if friendCheck {
		result.ResultCode = response.ErrFriendHasAlreadyExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("Friend has already exist, you don't need to request more")
	}

	// 3. Find exist friends request
	friendRequestEntityFromUserFound, err := user_entity.NewFriendRequest(command.FriendId, command.UserId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendRequestFromUserFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFromUserFound)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to check friend: %w", err)
	}

	if friendRequestFromUserFound {
		result.ResultCode = response.ErrFriendHasAlreadyExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("Your friend has already send add friend request, you don't need to request more")
	}

	friendRequestEntityFound, err := user_entity.NewFriendRequest(command.UserId, command.FriendId)

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFound)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to check friend request: %w", err)
	}

	// 4. Return if friend request has already exist
	if friendRequestFound {
		result.ResultCode = response.ErrFriendHasAlreadyExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("Friend request already exists, you don't need to request more")
	}

	// 5. Find user and friend
	userFound, err := s.userRepo.GetOne(ctx, "id=?", command.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("User record not found: %w", err)
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to get user: %w", err)
	}

	friendFound, err := s.userRepo.GetOne(ctx, "id=?", command.FriendId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("Friend record not found: %w", err)
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to get friend: %w", err)
	}

	// 6. Create friend request
	err = s.friendRequestRepo.CreateOne(ctx, &user_entity.FriendRequest{
		UserId:   command.UserId,
		FriendId: command.FriendId,
	})
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to add friend request: %w", err)
	}

	// 7. Push notification to user
	notificationEntity, err := notification_entity.NewNotification(
		userFound.FamilyName+" "+userFound.Name,
		userFound.AvatarUrl,
		friendFound.ID,
		consts.FRIEND_REQUEST,
		userFound.ID.String(),
		"",
	)

	_, err = s.notificationRepo.CreateOne(ctx, notificationEntity)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to add notification: %w", err)
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
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to send notification: %w", err)
	}

	// 9. Response success
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserFriend) GetFriendRequests(
	ctx context.Context,
	query *user_query.FriendRequestQuery,
) (result *user_query.FriendRequestQueryResult, err error) {
	result = &user_query.FriendRequestQueryResult{}
	// 1. Get list of user request to add friend
	userEntities, paging, err := s.friendRequestRepo.GetFriendRequests(ctx, query)
	if err != nil {
		result.Users = nil
		result.PagingResponse = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to get friend requests: %w", err)
	}

	// 2. Map userModel to userDtoShortVer
	var userResults []*common.UserShortVerResult
	for _, userEntity := range userEntities {
		userResults = append(userResults, mapper.NewUserShortVerEntity(userEntity))
	}

	result.Users = userResults
	result.PagingResponse = paging
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserFriend) AcceptFriendRequest(
	ctx context.Context,
	command *user_command.AcceptFriendRequestCommand,
) (result *user_command.AcceptFriendRequestCommandResult, err error) {
	result = &user_command.AcceptFriendRequestCommandResult{}
	// 1. Find exist friends request
	friendRequestEntityFound, err := user_entity.NewFriendRequest(command.UserId, command.FriendId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFound)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to check friend request: %w", err)
	}

	// 2. Return if friend request is not exist
	if !friendRequestFound {
		result.ResultCode = response.ErrFriendNotExist
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("Friend request is not exist: %w", err)
	}

	// 3. Find user and friend
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("User record not found: %w", err)
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to get user: %w", err)
	}

	friendFound, err := s.userRepo.GetById(ctx, command.FriendId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("Friend record not found: %w", err)
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to get friend: %w", err)
	}

	// 4. Create friend
	friendEntityForUser, err := user_entity.NewFriend(userFound.ID, friendFound.ID)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendEntityForFriend, err := user_entity.NewFriend(friendFound.ID, userFound.ID)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	err = s.friendRepo.CreateOne(ctx, friendEntityForUser)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to add friend: %w", err)
	}

	err = s.friendRepo.CreateOne(ctx, friendEntityForFriend)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to add friend: %w", err)
	}

	// 5. Delete friendRequest
	err = s.friendRequestRepo.DeleteOne(ctx, &user_entity.FriendRequest{
		UserId:   command.UserId,
		FriendId: command.FriendId,
	})
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to delete friend request: %w", err)
	}

	// 6. Push notification to user
	notificationModel := &notification_entity.Notification{
		From:             friendFound.FamilyName + " " + friendFound.Name,
		FromUrl:          friendFound.AvatarUrl,
		UserId:           userFound.ID,
		NotificationType: consts.ACCEPT_FRIEND_REQUEST,
		ContentId:        (friendFound.ID).String(),
	}

	_, err = s.notificationRepo.CreateOne(ctx, notificationModel)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to add notification: %w", err)
	}

	// 7. Send realtime notification (websocket)
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
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to send notification: %w", err)
	}

	// 8. Response success
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserFriend) RemoveFriendRequest(
	ctx context.Context,
	command *user_command.RemoveFriendRequestCommand,
) (result *user_command.RemoveFriendRequestCommandResult, err error) {
	result = &user_command.RemoveFriendRequestCommandResult{}
	// 1. Find exist friends request
	friendRequestEntityFound, err := user_entity.NewFriendRequest(command.UserId, command.FriendId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFound)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to check friend request: %w", err)
	}

	// 2. Return if friend request is not exist
	if !friendRequestFound {
		result.ResultCode = response.ErrFriendRequestNotExists
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("Friend request is not exist: %w", err)
	}

	// 3. Delete friend request
	friendRequestEntity, err := user_entity.NewFriendRequest(command.UserId, command.FriendId)

	err = s.friendRequestRepo.DeleteOne(ctx, friendRequestEntity)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to delete friend request: %w", err)
	}

	// 4. Response success
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusNoContent
	return result, nil
}

func (s *sUserFriend) UnFriend(
	ctx context.Context,
	command *user_command.UnFriendCommand,
) (result *user_command.UnFriendCommandResult, err error) {
	result = &user_command.UnFriendCommandResult{}
	// 1. Check friend exist
	friendEntity, err := user_entity.NewFriend(command.UserId, command.FriendId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	friendCheck, err := s.friendRepo.CheckFriendExist(ctx, friendEntity)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to check friend: %w", err)
	}

	if !friendCheck {
		result.ResultCode = response.ErrFriendNotExist
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("Friend is not exist: %w", err)
	}

	// 2. Remove friend
	err = s.friendRepo.DeleteOne(ctx, friendEntity)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to delete friend: %w", err)
	}

	friendEntityForFriend, err := user_entity.NewFriend(command.FriendId, command.UserId)
	if err != nil {
		result.ResultCode = response.ErrCodeValidate
		result.HttpStatusCode = http.StatusBadRequest
		return result, err
	}

	err = s.friendRepo.DeleteOne(ctx, friendEntityForFriend)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to delete friend: %w", err)
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserFriend) GetFriends(
	ctx context.Context,
	query *user_query.FriendQuery,
) (result *user_query.FriendQueryResult, err error) {
	result = &user_query.FriendQueryResult{}
	// 1. Get list of friend
	userEntities, paging, err := s.friendRepo.GetFriends(ctx, query)
	if err != nil {
		result.Users = nil
		result.PagingResponse = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Failed to get friends: %w", err)
	}

	// 2. Map userModel to userResultShortVer
	var userResults []*common.UserShortVerResult
	for _, userEntity := range userEntities {
		userResults = append(userResults, mapper.NewUserShortVerEntity(userEntity))
	}

	result.Users = userResults
	result.PagingResponse = paging
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
