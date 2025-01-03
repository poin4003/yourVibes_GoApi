package implement

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	userMapper "github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	userQuery "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	userRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/media"
	"gorm.io/gorm"
)

type sUserInfo struct {
	userRepo          userRepo.IUserRepository
	settingRepo       userRepo.ISettingRepository
	friendRepo        userRepo.IFriendRepository
	friendRequestRepo userRepo.IFriendRequestRepository
}

func NewUserInfoImplement(
	userRepo userRepo.IUserRepository,
	settingRepo userRepo.ISettingRepository,
	friendRepo userRepo.IFriendRepository,
	friendRequestRepo userRepo.IFriendRequestRepository,
) *sUserInfo {
	return &sUserInfo{
		userRepo:          userRepo,
		settingRepo:       settingRepo,
		friendRepo:        friendRepo,
		friendRequestRepo: friendRequestRepo,
	}
}

func (s *sUserInfo) GetInfoByUserId(
	ctx context.Context,
	query *userQuery.GetOneUserQuery,
) (result *userQuery.UserQueryResult, err error) {
	result = &userQuery.UserQueryResult{
		User:           nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Find User
	userFound, err := s.userRepo.GetOne(ctx, "id = ?", query.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.User = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Return if user fetches his own information
	if query.AuthenticatedUserId == query.UserId {
		result.User = userMapper.NewUserResultWithoutSettingEntity(userFound, consts.NOT_FRIEND)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}

	// 3. Check friend status
	var friendStatus consts.FriendStatus
	isFriend, err := s.friendRepo.CheckFriendExist(ctx, &userEntity.Friend{
		UserId:   query.AuthenticatedUserId,
		FriendId: query.UserId,
	})
	if err != nil {
		return result, err
	}

	// 3.1. Check friend
	if isFriend {
		friendStatus = consts.IS_FRIEND
	} else {
		// 3.2. Check if user are send add friend request
		isSendFriendRequest, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, &userEntity.FriendRequest{
			UserId:   query.AuthenticatedUserId,
			FriendId: query.UserId,
		})
		if err != nil {
			return result, err
		}
		if isSendFriendRequest {
			friendStatus = consts.SEND_FRIEND_REQUEST
		} else {
			// 3.3. Check if user are receive add friend request
			isReceiveFriendRequest, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, &userEntity.FriendRequest{
				UserId:   query.UserId,
				FriendId: query.AuthenticatedUserId,
			})
			if err != nil {
				return result, err
			}
			if isReceiveFriendRequest {
				friendStatus = consts.RECEIVE_FRIEND_REQUEST
			} else {
				friendStatus = consts.NOT_FRIEND
			}
		}
	}

	// 4. Check privacy
	var resultCode int
	var userResult *common.UserWithoutSettingResult
	switch userFound.Privacy {
	case consts.PUBLIC:
		userResult = userMapper.NewUserResultWithoutSettingEntity(userFound, friendStatus)
		resultCode = response.ErrCodeSuccess
	case consts.FRIEND_ONLY:
		if friendStatus == consts.IS_FRIEND {
			userResult = userMapper.NewUserResultWithoutSettingEntity(userFound, friendStatus)
			resultCode = response.ErrCodeSuccess
		} else {
			userResult = userMapper.NewUserResultWithoutPrivateInfo(userFound, friendStatus)
			resultCode = response.ErrUserFriendAccess
		}
	case consts.PRIVATE:
		userResult = userMapper.NewUserResultWithoutPrivateInfo(userFound, friendStatus)
		resultCode = response.ErrUserPrivateAccess
	default:
		userResult = userMapper.NewUserResultWithoutPrivateInfo(userFound, friendStatus)
		resultCode = response.ErrUserPrivateAccess
	}

	result.User = userResult
	result.ResultCode = resultCode
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserInfo) GetManyUsers(
	ctx context.Context,
	query *userQuery.GetManyUserQuery,
) (result *userQuery.UserQueryListResult, err error) {
	result = &userQuery.UserQueryListResult{
		Users:          nil,
		PagingResponse: nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	userEntities, paging, err := s.userRepo.GetMany(ctx, query)

	if err != nil {
		return result, err
	}

	var userResultList []*common.UserShortVerResult
	for _, user := range userEntities {
		userResultList = append(userResultList, userMapper.NewUserShortVerEntity(user))
	}

	result.Users = userResultList
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}

func (s *sUserInfo) UpdateUser(
	ctx context.Context,
	command *userCommand.UpdateUserCommand,
) (result *userCommand.UpdateUserCommandResult, err error) {
	result = &userCommand.UpdateUserCommandResult{
		User:           nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. update setting language
	if command.LanguageSetting != nil {
		settingFound, err := s.settingRepo.GetSetting(ctx, "user_id=?", command.UserId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			return result, fmt.Errorf("failed to get setting for user %v: %w", command.UserId, err)
		}
		_, err = s.settingRepo.UpdateOne(ctx, settingFound.ID,
			&userEntity.SettingUpdate{Language: command.LanguageSetting},
		)
	}

	// 2. update user information
	updateUserEntity := &userEntity.UserUpdate{
		FamilyName:  command.FamilyName,
		Name:        command.Name,
		PhoneNumber: command.PhoneNumber,
		Birthday:    command.Birthday,
		Privacy:     command.Privacy,
		Biography:   command.Biography,
	}

	err = updateUserEntity.ValidateUserUpdate()
	if err != nil {
		return result, err
	}

	// 3. update Avatar
	if command.Avatar != nil && command.Avatar.Size > 0 && command.Avatar.Filename != "" {
		avatarUrl, err := media.SaveMedia(command.Avatar)
		if err != nil {
			return result, fmt.Errorf("failed to upload Avatar: %w", err)
		}

		_, err = s.userRepo.UpdateOne(ctx, *command.UserId, &userEntity.UserUpdate{
			AvatarUrl: &avatarUrl,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			return result, err
		}
	}

	// 4. update Capwall
	if command.Capwall != nil && command.Capwall.Size > 0 && command.Capwall.Filename != "" {
		capwallUrl, err := media.SaveMedia(command.Capwall)
		if err != nil {
			return result, fmt.Errorf("failed to upload Capwall: %w", err)
		}

		_, err = s.userRepo.UpdateOne(ctx, *command.UserId, &userEntity.UserUpdate{
			CapwallUrl: &capwallUrl,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			return result, err
		}
	}

	userFound, err := s.userRepo.UpdateOne(ctx, *command.UserId, updateUserEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	result.User = userMapper.NewUserResultFromEntity(userFound)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserInfo) GetUserStatusById(
	ctx context.Context,
	id uuid.UUID,
) (status bool, err error) {
	userStatus, err := s.userRepo.GetStatusById(ctx, id)
	if err != nil {
		return false, err
	}
	return userStatus, err
}
