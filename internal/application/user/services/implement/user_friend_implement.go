package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	user_mapper "github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	user_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/mapper"
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
	// 1. Check exist friend
	friendCheck, err := s.friendRepo.CheckFriendExist(ctx, &user_entity.Friend{
		UserId:   command.UserId,
		FriendId: command.FriendId,
	})
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
	friendRequestFromUserFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, &user_entity.FriendRequest{
		UserId:   command.FriendId,
		FriendId: command.UserId,
	})
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

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, &user_entity.FriendRequest{
		UserId:   command.UserId,
		FriendId: command.FriendId,
	})
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
	notificationEntity, err := user_entity.NewNotification(
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
	for i, userEntity := range userEntities {
		userResult := user_mapper.NewUserShortVerEntity(userEntity)
		result.Users[i] = *userResult
	}

	result.PagingResponse = paging
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserFriend) AcceptFriendRequest(
	ctx context.Context,
	command *user_command.AcceptFriendRequestCommand,
) (result *user_command.AcceptFriendRequestCommandResult, err error) {
	// 1. Find exist friends request
	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to check friend request: %w", err)
	}

	// 2. Return if friend request is not exist
	if !friendRequestFound {
		return response.ErrFriendRequestNotExists, http.StatusBadRequest, fmt.Errorf("Friend request is not exist: %w", err)
	}

	// 3. Find user and friend
	userFound, err := s.userRepo.GetUser(ctx, "id=?", friendRequest.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrDataNotFound, http.StatusBadRequest, fmt.Errorf("User record not found: %w", err)
		}
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to get user: %w", err)
	}

	friendFound, err := s.userRepo.GetUser(ctx, "id=?", friendRequest.FriendId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrDataNotFound, http.StatusBadRequest, fmt.Errorf("Friend record not found: %w", err)
		}
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to get friend: %w", err)
	}

	// 4. Create friend
	friendModelForUser := &entities2.Friend{
		UserId:   userFound.ID,
		FriendId: friendFound.ID,
	}

	friendModelForFriend := &entities2.Friend{
		UserId:   friendFound.ID,
		FriendId: userFound.ID,
	}

	err = s.friendRepo.CreateFriend(ctx, friendModelForUser)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to add friend: %w", err)
	}

	err = s.friendRepo.CreateFriend(ctx, friendModelForFriend)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to add friend: %w", err)
	}

	// 5. Delete friendRequest
	err = s.friendRequestRepo.DeleteFriendRequest(ctx, friendRequest)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to delete friend request: %w", err)
	}

	// 6. Push notification to user
	notificationModel := &entities2.Notification{
		From:             friendFound.FamilyName + " " + friendFound.Name,
		FromUrl:          friendFound.AvatarUrl,
		UserId:           userFound.ID,
		NotificationType: consts.ACCEPT_FRIEND_REQUEST,
		ContentId:        (friendFound.ID).String(),
	}
	notification, err := s.notificationRepo.CreateNotification(ctx, notificationModel)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to add notification: %w", err)
	}

	// 7. Send realtime notification (websocket)
	notificationDto := mapper.MapNotificationToNotificationDto(notification)

	err = global.SocketHub.SendNotification(userFound.ID.String(), notificationDto)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to send notification: %w", err)
	}

	// 8. Response success
	return response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sUserFriend) RemoveFriendRequest(
	ctx context.Context,
	command *user_command.RemoveFriendRequestCommand,
) (result *user_command.RemoveFriendRequestCommandResult, err error) {
	// 1. Find exist friends request
	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequest)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to check friend request: %w", err)
	}

	// 2. Return if friend request is not exist
	if !friendRequestFound {
		return response.ErrFriendRequestNotExists, http.StatusBadRequest, fmt.Errorf("Friend request is not exist: %w", err)
	}

	// 3. Delete friend request
	err = s.friendRequestRepo.DeleteFriendRequest(ctx, friendRequest)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to delete friend request: %w", err)
	}

	// 4. Response success
	return response.ErrCodeSuccess, http.StatusNoContent, nil
}

func (s *sUserFriend) UnFriend(
	ctx context.Context,
	command *user_command.UnFriendCommand,
) (result *user_command.UnFriendCommandResult, err error) {
	// 1. Check friend exist
	friendCheck, err := s.friendRepo.CheckFriendExist(ctx, friend)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to check friend: %w", err)
	}

	if !friendCheck {
		return response.ErrFriendNotExist, http.StatusBadRequest, fmt.Errorf("Friend is not exist: %w", err)
	}

	// 2. Remove friend
	err = s.friendRepo.DeleteFriend(ctx, friend)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to delete friend: %w", err)
	}

	friendModelForFriend := &entities2.Friend{
		UserId:   friend.FriendId,
		FriendId: friend.UserId,
	}

	err = s.friendRepo.DeleteFriend(ctx, friendModelForFriend)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to delete friend: %w", err)
	}

	return response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sUserFriend) GetFriends(
	ctx context.Context,
	query *user_query.FriendQuery,
) (result *user_query.FriendQueryResult, err error) {
	// 1. Get list of friend
	userModels, paging, err := s.friendRepo.GetFriends(ctx, userId, query)
	if err != nil {
		return nil, nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to get friends: %w", err)
	}

	// 2. Map userModel to userDtoShortVer
	for _, userModel := range userModels {
		userDto := mapper.MapUserToUserDtoShortVer(userModel)
		userDtos = append(userDtos, &userDto)
	}

	return userDtos, paging, response.ErrCodeSuccess, http.StatusOK, nil
}
