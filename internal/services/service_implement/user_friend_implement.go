package service_implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/user_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
)

type sUserFriend struct {
	userRepo          repository.IUserRepository
	friendRequestRepo repository.IFriendRequestRepository
	friendRepo        repository.IFriendRepository
	notificationRepo  repository.INotificationRepository
}

func NewUserFriendImplement(
	userRepo repository.IUserRepository,
	friendRequestRepo repository.IFriendRequestRepository,
	friendRepo repository.IFriendRepository,
	notificationRepo repository.INotificationRepository,
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
	friendRequest *model.FriendRequest,
) (resultCode int, httpStatusCode int, err error) {
	// 1. Check exist friend
	friendModel := &model.Friend{
		UserId:   friendRequest.UserId,
		FriendId: friendRequest.FriendId,
	}
	friendCheck, err := s.friendRepo.CheckFriendExist(ctx, friendModel)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to check friend: %w", err)
	}

	// 2. Return if friend has already exist
	if friendCheck {
		return response.ErrFriendHasAlreadyExists, http.StatusBadRequest, fmt.Errorf("Friend has already exist, you don't need to request more")
	}

	// 3. Find exist friends request
	friendRequestFromUserModel := &model.FriendRequest{
		UserId:   friendRequest.FriendId,
		FriendId: friendRequest.UserId,
	}

	friendRequestFromUserFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestFromUserModel)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to check friend: %w", err)
	}

	if friendRequestFromUserFound {
		return response.ErrFriendRequestHasAlreadyExists, http.StatusBadRequest, fmt.Errorf("Your friend has already send add friend request, you don't need to request more")
	}

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequest)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to check friend request: %w", err)
	}

	// 4. Return if friend request has already exist
	if friendRequestFound {
		return response.ErrFriendRequestHasAlreadyExists, http.StatusBadRequest, fmt.Errorf("Friend request already exists, you don't need to request more")
	}

	// 5. Find user and friend
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

	// 6. Create friend request
	err = s.friendRequestRepo.CreateFriendRequest(ctx, friendRequest)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to add friend request: %w", err)
	}

	// 7. Push notification to user
	notificationModel := &model.Notification{
		From:             userFound.FamilyName + " " + userFound.Name,
		FromUrl:          userFound.AvatarUrl,
		UserId:           friendFound.ID,
		NotificationType: consts.FRIEND_REQUEST,
		ContentId:        (userFound.ID).String(),
	}
	notification, err := s.notificationRepo.CreateNotification(ctx, notificationModel)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to add notification: %w", err)
	}

	// 8. Send realtime notification (websocket)
	notificationDto := mapper.MapNotificationToNotificationDto(notification)

	err = global.SocketHub.SendNotification(friendFound.ID.String(), notificationDto)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to send notification: %w", err)
	}

	// 9. Response success
	return response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sUserFriend) GetFriendRequests(
	ctx context.Context,
	userId uuid.UUID,
	query *query_object.FriendRequestQueryObject,
) (userDtos []*user_dto.UserDtoShortVer, pagingResponse *response.PagingResponse, resultCode int, httpStatusCode int, err error) {
	// 1. Get list of user request to add friend
	userModels, paging, err := s.friendRequestRepo.GetFriendRequest(ctx, userId, query)
	if err != nil {
		return nil, nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Failed to get friend requests: %w", err)
	}

	// 2. Map userModel to userDtoShortVer
	for _, userModel := range userModels {
		userDto := mapper.MapUserToUserDtoShortVer(userModel)
		userDtos = append(userDtos, &userDto)
	}

	return userDtos, paging, response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sUserFriend) AcceptFriendRequest(
	ctx context.Context,
	friendRequest *model.FriendRequest,
) (resultCode int, httpStatusCode int, err error) {
	// 1. Find exist friends request
	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequest)
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
	friendModelForUser := &model.Friend{
		UserId:   userFound.ID,
		FriendId: friendFound.ID,
	}

	friendModelForFriend := &model.Friend{
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
	notificationModel := &model.Notification{
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
	friendRequest *model.FriendRequest,
) (resultCode int, httpStatusCode int, err error) {
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
	friend *model.Friend,
) (resultCode int, httpStatusCode int, err error) {
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

	friendModelForFriend := &model.Friend{
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
	userId uuid.UUID,
	query *query_object.FriendQueryObject,
) (userDtos []*user_dto.UserDtoShortVer, pagingResponse *response.PagingResponse, resultCode int, httpStatusCode int, err error) {
	// 1. Get list of friend
	userModels, paging, err := s.friendRepo.GetFriend(ctx, userId, query)
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
