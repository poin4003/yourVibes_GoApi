package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/global"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/producer"
	userQuery "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"
	"go.uber.org/zap"
	"sort"
	"sync"
)

type sUserFriend struct {
	userRepo              repositories.IUserRepository
	friendRequestRepo     repositories.IFriendRequestRepository
	friendRepo            repositories.IFriendRepository
	userCache             cache.IUserCache
	notificationPublisher *producer.NotificationPublisher
}

func NewUserFriendImplement(
	userRepo repositories.IUserRepository,
	friendRequestRepo repositories.IFriendRequestRepository,
	friendRepo repositories.IFriendRepository,
	userCache cache.IUserCache,
	notificationPublisher *producer.NotificationPublisher,
) *sUserFriend {
	return &sUserFriend{
		userRepo:              userRepo,
		friendRequestRepo:     friendRequestRepo,
		friendRepo:            friendRepo,
		userCache:             userCache,
		notificationPublisher: notificationPublisher,
	}
}

func (s *sUserFriend) SendAddFriendRequest(
	ctx context.Context,
	command *userCommand.SendAddFriendRequestCommand,
) (err error) {
	// 1. Check exist friend
	friendEntity, err := userEntity.NewFriend(command.UserId, command.FriendId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	friendCheck, err := s.friendRepo.CheckFriendExist(ctx, friendEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 2. Return if friend has already exist
	if friendCheck {
		return response.NewCustomError(
			response.ErrFriendHasAlreadyExists,
			"friend has already exist, you don't need to request more",
		)
	}

	// 3. Find exist friends request
	friendRequestEntityFromUserFound, err := userEntity.NewFriendRequest(command.FriendId, command.UserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	friendRequestFromUserFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFromUserFound)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if friendRequestFromUserFound {
		return response.NewCustomError(
			response.ErrFriendHasAlreadyExists,
			"your friend has already send add friend request, you don't need to request more",
		)
	}

	friendRequestEntityFound, _ := userEntity.NewFriendRequest(command.UserId, command.FriendId)

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFound)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Return if friend request has already exist
	if friendRequestFound {
		return response.NewCustomError(
			response.ErrFriendHasAlreadyExists,
			"friend request already exists, you don't need to request more",
		)
	}

	// 5. Find user and friend
	userFound, err := s.userRepo.GetOne(ctx, "id=?", command.UserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return response.NewCustomError(response.UserNotFound)
	}

	friendFound, err := s.userRepo.GetOne(ctx, "id=?", command.FriendId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if friendFound == nil {
		return response.NewCustomError(response.UserNotFound)
	}

	// 6. Create friend request
	err = s.friendRequestRepo.CreateOne(ctx, &userEntity.FriendRequest{
		UserId:   command.UserId,
		FriendId: command.FriendId,
	})
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 7. Publish to RabbitMQ to handle notification
	notification, _ := notificationEntity.NewNotification(
		userFound.FamilyName+" "+userFound.Name,
		userFound.AvatarUrl,
		friendFound.ID,
		consts.FRIEND_REQUEST,
		userFound.ID.String(),
		"",
	)

	notiMsg := mapper.NewNotificationResult(notification)
	if err = s.notificationPublisher.PublishNotification(ctx, notiMsg, "notification.single.db_websocket"); err != nil {
		global.Logger.Error("Failed to publish notification", zap.Error(err))
	}

	// 9. Response success
	return nil
}

func (s *sUserFriend) GetFriendRequests(
	ctx context.Context,
	query *userQuery.FriendRequestQuery,
) (result *userQuery.FriendRequestQueryResult, err error) {
	// 1. Get list of user request to add friend
	userEntities, paging, err := s.friendRequestRepo.GetFriendRequests(ctx, query)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 2. Map userEntity to userDtoShortVer
	var userResults []*common.UserShortVerResult
	for _, user := range userEntities {
		userResults = append(userResults, mapper.NewUserShortVerEntity(user))
	}

	return &userQuery.FriendRequestQueryResult{
		Users:          userResults,
		PagingResponse: paging,
	}, nil
}

func (s *sUserFriend) AcceptFriendRequest(
	ctx context.Context,
	command *userCommand.AcceptFriendRequestCommand,
) (err error) {
	// 1. Find exist friends request
	friendRequestEntityFound, err := userEntity.NewFriendRequest(command.UserId, command.FriendId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFound)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 2. Return if friend request is not exist
	if !friendRequestFound {
		return response.NewCustomError(response.ErrFriendNotExist)
	}

	// 3. Find user and friend
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return response.NewDataNotFoundError("user not found")
	}

	friendFound, err := s.userRepo.GetById(ctx, command.FriendId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if friendFound == nil {
		return response.NewDataNotFoundError("friend not found")
	}

	// 4. Create friend
	friendEntityForUser, err := userEntity.NewFriend(userFound.ID, friendFound.ID)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	friendEntityForFriend, err := userEntity.NewFriend(friendFound.ID, userFound.ID)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	err = s.friendRepo.CreateOne(ctx, friendEntityForUser)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	err = s.friendRepo.CreateOne(ctx, friendEntityForFriend)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 5. Delete friendRequest
	err = s.friendRequestRepo.DeleteOne(ctx, &userEntity.FriendRequest{
		UserId:   command.UserId,
		FriendId: command.FriendId,
	})
	if err != nil {
		return response.NewServerFailedError(err.Error())
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
		return response.NewServerFailedError(err.Error())
	}

	err = updateFriendData.ValidateUserUpdate()
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, updateUserData)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.userRepo.UpdateOne(ctx, friendFound.ID, updateFriendData)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 7. Push notification to user
	notification := &notificationEntity.Notification{
		From:             friendFound.FamilyName + " " + friendFound.Name,
		FromUrl:          friendFound.AvatarUrl,
		UserId:           userFound.ID,
		NotificationType: consts.ACCEPT_FRIEND_REQUEST,
		ContentId:        (friendFound.ID).String(),
	}

	notiMsg := mapper.NewNotificationResult(notification)
	if err = s.notificationPublisher.PublishNotification(ctx, notiMsg, "notification.single.db_websocket"); err != nil {
		global.Logger.Error("Failed to publish notification", zap.Error(err))
	}

	// 9. Response success
	return nil
}

func (s *sUserFriend) RemoveFriendRequest(
	ctx context.Context,
	command *userCommand.RemoveFriendRequestCommand,
) (err error) {
	// 1. Find exist friends request
	friendRequestEntityFound, err := userEntity.NewFriendRequest(command.UserId, command.FriendId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	friendRequestFound, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, friendRequestEntityFound)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 2. Return if friend request is not exist
	if !friendRequestFound {
		return response.NewCustomError(response.ErrFriendRequestNotExists)
	}

	// 3. Delete friend request
	friendRequestEntity, _ := userEntity.NewFriendRequest(command.UserId, command.FriendId)

	err = s.friendRequestRepo.DeleteOne(ctx, friendRequestEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Response success
	return nil
}

func (s *sUserFriend) UnFriend(
	ctx context.Context,
	command *userCommand.UnFriendCommand,
) (err error) {
	// 1. Check friend exist
	friendEntity, err := userEntity.NewFriend(command.UserId, command.FriendId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	friendCheck, err := s.friendRepo.CheckFriendExist(ctx, friendEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if !friendCheck {
		return response.NewCustomError(response.ErrFriendNotExist)
	}

	// 2. Remove friend
	err = s.friendRepo.DeleteOne(ctx, friendEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	friendEntityForFriend, err := userEntity.NewFriend(command.FriendId, command.UserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	err = s.friendRepo.DeleteOne(ctx, friendEntityForFriend)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 3. Minus -1 to friend count of user and friend
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return response.NewDataNotFoundError("can not found user")
	}

	friendFound, err := s.userRepo.GetById(ctx, friendEntityForFriend.UserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if friendFound == nil {
		return response.NewDataNotFoundError("can not found friend")
	}

	updateUserData := &userEntity.UserUpdate{
		FriendCount: pointer.Ptr(userFound.FriendCount - 1),
	}

	updateFriendData := &userEntity.UserUpdate{
		FriendCount: pointer.Ptr(friendFound.FriendCount - 1),
	}

	err = updateUserData.ValidateUserUpdate()
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	err = updateFriendData.ValidateUserUpdate()
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, updateUserData)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.userRepo.UpdateOne(ctx, friendFound.ID, updateFriendData)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sUserFriend) GetFriends(
	ctx context.Context,
	query *userQuery.FriendQuery,
) (result *userQuery.FriendQueryResult, err error) {
	// 1. Get list of friend
	userEntities, paging, err := s.friendRepo.GetFriends(ctx, query)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 2. Map userModel to userResultShortVerWithActiveStatus
	var wg sync.WaitGroup
	userResultsChan := make(chan *common.UserShortVerWithActiveStatusResult, len(userEntities))
	for _, user := range userEntities {
		wg.Add(1)
		go func(user *userEntity.User) {
			defer wg.Done()
			userActiveStatus := s.userCache.IsOnline(ctx, user.ID)
			userResultsChan <- mapper.NewUserShortVerWithActiveStatus(user, userActiveStatus)
		}(user)
	}

	wg.Wait()
	close(userResultsChan)

	var userResults []*common.UserShortVerWithActiveStatusResult
	for userResult := range userResultsChan {
		userResults = append(userResults, userResult)
	}

	sort.Slice(userResults, func(i, j int) bool {
		return userResults[i].UserActiveStatus && !userResults[j].UserActiveStatus
	})

	return &userQuery.FriendQueryResult{
		Users:          userResults,
		PagingResponse: paging,
	}, nil
}

func (s *sUserFriend) GetFriendSuggestion(
	ctx context.Context,
	query *userQuery.FriendQuery,
) (result *userQuery.FriendWithIsSendFriendRequestQueryResult, err error) {
	// 1. Get list friend suggestion
	userEntities, paging, err := s.friendRepo.GetFriendSuggestions(ctx, query)
	if err != nil {
		return nil, err
	}

	// 2. Map to return
	var userResults []*common.UserShortVerWithSendFriendRequestResult
	for _, user := range userEntities {
		userResults = append(userResults, mapper.NewUserShortVerWithSendFriendRequestEntity(user))
	}

	return &userQuery.FriendWithIsSendFriendRequestQueryResult{
		Users:          userResults,
		PagingResponse: paging,
	}, nil
}

func (s *sUserFriend) GetFriendByBirthday(
	ctx context.Context,
	query *userQuery.FriendQuery,
) (result *userQuery.FriendWithBirthdayQueryResult, err error) {
	// 1. Get list friend with birthday
	userEntities, paging, err := s.friendRepo.GetFriendByBirthday(ctx, query)
	if err != nil {
		return nil, err
	}

	// 2. Map to return
	var userResults []*common.UserShortVerWithBirthdayResult
	for _, user := range userEntities {
		userResults = append(userResults, mapper.NewUserShortVerWithBirthdayEntity(user))
	}

	return &userQuery.FriendWithBirthdayQueryResult{
		Users:          userResults,
		PagingResponse: paging,
	}, nil
}

func (s *sUserFriend) GetNonFriends(
	ctx context.Context,
	query *userQuery.FriendQuery,
) (result *userQuery.FriendWithIsSendFriendRequestQueryResult, err error) {
	// 1. Get list non friend
	userEntities, paging, err := s.friendRepo.GetNonFriends(ctx, query)
	if err != nil {
		return nil, err
	}

	// 2. Map to return
	var userResults []*common.UserShortVerWithSendFriendRequestResult
	for _, user := range userEntities {
		userResults = append(userResults, mapper.NewUserShortVerWithSendFriendRequestEntity(user))
	}

	return &userQuery.FriendWithIsSendFriendRequestQueryResult{
		Users:          userResults,
		PagingResponse: paging,
	}, nil
}
