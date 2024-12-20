package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	user_mapper "github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	user_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/cloudinary_util"
	"gorm.io/gorm"
	"net/http"
)

type sUserInfo struct {
	userRepo          user_repo.IUserRepository
	settingRepo       user_repo.ISettingRepository
	friendRepo        user_repo.IFriendRepository
	friendRequestRepo user_repo.IFriendRequestRepository
}

func NewUserInfoImplement(
	userRepo user_repo.IUserRepository,
	settingRepo user_repo.ISettingRepository,
	friendRepo user_repo.IFriendRepository,
	friendRequestRepo user_repo.IFriendRequestRepository,
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
	query *user_query.GetOneUserQuery,
) (result *user_query.UserQueryResult, err error) {
	result = &user_query.UserQueryResult{}
	// 1. Find User
	userFound, err := s.userRepo.GetOne(ctx, "id = ?", query.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.User = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.User = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 2. Return if user fetches his own information
	if query.AuthenticatedUserId == query.UserId {
		result.User = user_mapper.NewUserResultWithoutSettingEntity(userFound, consts.NOT_FRIEND)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}

	// 3. Check friend status
	var friendStatus consts.FriendStatus
	isFriend, err := s.friendRepo.CheckFriendExist(ctx, &user_entity.Friend{
		UserId:   query.AuthenticatedUserId,
		FriendId: query.UserId,
	})
	if err != nil {
		result.User = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 3.1. Check friend
	if isFriend {
		friendStatus = consts.IS_FRIEND
	} else {
		// 3.2. Check if user are send add friend request
		isSendFriendRequest, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, &user_entity.FriendRequest{
			UserId:   query.AuthenticatedUserId,
			FriendId: query.UserId,
		})
		if err != nil {
			result.User = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, err
		}
		if isSendFriendRequest {
			friendStatus = consts.SEND_FRIEND_REQUEST
		} else {
			// 3.3. Check if user are receive add friend request
			isReceiveFriendRequest, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, &user_entity.FriendRequest{
				UserId:   query.UserId,
				FriendId: query.AuthenticatedUserId,
			})
			if err != nil {
				result.User = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
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
		userResult = user_mapper.NewUserResultWithoutSettingEntity(userFound, friendStatus)
		resultCode = response.ErrCodeSuccess
	case consts.FRIEND_ONLY:
		if friendStatus == consts.IS_FRIEND {
			userResult = user_mapper.NewUserResultWithoutSettingEntity(userFound, friendStatus)
			resultCode = response.ErrCodeSuccess
		} else {
			userResult = user_mapper.NewUserResultWithoutPrivateInfo(userFound, friendStatus)
			resultCode = response.ErrUserFriendAccess
		}
	case consts.PRIVATE:
		userResult = user_mapper.NewUserResultWithoutPrivateInfo(userFound, friendStatus)
		resultCode = response.ErrUserPrivateAccess
	default:
		userResult = user_mapper.NewUserResultWithoutPrivateInfo(userFound, friendStatus)
		resultCode = response.ErrUserPrivateAccess
	}

	result.User = userResult
	result.ResultCode = resultCode
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserInfo) GetManyUsers(
	ctx context.Context,
	query *user_query.GetManyUserQuery,
) (result *user_query.UserQueryListResult, err error) {
	result = &user_query.UserQueryListResult{}
	userEntities, paging, err := s.userRepo.GetMany(ctx, query)

	if err != nil {
		result.Users = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.PagingResponse = nil
		return result, err
	}

	var userResultList []*common.UserShortVerResult
	for _, userEntity := range userEntities {
		userResultList = append(userResultList, user_mapper.NewUserShortVerEntity(userEntity))
	}

	result.Users = userResultList
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}

func (s *sUserInfo) UpdateUser(
	ctx context.Context,
	command *user_command.UpdateUserCommand,
) (result *user_command.UpdateUserCommandResult, err error) {
	result = &user_command.UpdateUserCommandResult{}
	// 1. update setting language
	if command.LanguageSetting != nil {
		settingFound, err := s.settingRepo.GetSetting(ctx, "user_id=?", command.UserId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.User = nil
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			result.User = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("Failed to get setting for user %v: %w", command.UserId, err)
		}
		_, err = s.settingRepo.UpdateOne(ctx, settingFound.ID,
			&user_entity.SettingUpdate{Language: command.LanguageSetting},
		)
	}

	// 2. update user information
	updateUserEntity := &user_entity.UserUpdate{
		FamilyName:  command.FamilyName,
		Name:        command.Name,
		PhoneNumber: command.PhoneNumber,
		Birthday:    command.Birthday,
		Privacy:     command.Privacy,
		Biography:   command.Biography,
	}

	err = updateUserEntity.ValidateUserUpdate()
	if err != nil {
		result.User = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 3. update Avatar
	if command.Avatar != nil {
		avatarUrl, err := cloudinary_util.UploadMediaToCloudinary(command.Avatar)
		if err != nil {
			result.User = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to upload Avatar: %w", err)
		}

		_, err = s.userRepo.UpdateOne(ctx, *command.UserId, &user_entity.UserUpdate{
			AvatarUrl: &avatarUrl,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.User = nil
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			result.User = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, err
		}
	}

	// 4. update Capwall
	if command.Capwall != nil {
		capwallUrl, err := cloudinary_util.UploadMediaToCloudinary(command.Capwall)
		if err != nil {
			result.User = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to upload Capwall: %w", err)
		}

		_, err = s.userRepo.UpdateOne(ctx, *command.UserId, &user_entity.UserUpdate{
			CapwallUrl: &capwallUrl,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.User = nil
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			result.User = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, err
		}
	}

	userFound, err := s.userRepo.UpdateOne(ctx, *command.UserId, updateUserEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.User = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.User = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	result.User = user_mapper.NewUserResultFromEntity(userFound)
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
