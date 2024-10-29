package service_implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/user_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/cloudinary_util"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"

	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
)

type sUserInfo struct {
	userRepo          repository.IUserRepository
	settingRepo       repository.ISettingRepository
	friendRepo        repository.IFriendRepository
	friendRequestRepo repository.IFriendRequestRepository
}

func NewUserInfoImplement(
	userRepo repository.IUserRepository,
	settingRepo repository.ISettingRepository,
	friendRepo repository.IFriendRepository,
	friendRequestRepo repository.IFriendRequestRepository,
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
	userId uuid.UUID,
	authenticatedUserId uuid.UUID,
) (userDto *user_dto.UserDtoWithoutSetting, resultCode int, httpStatusCode int, err error) {
	// 1. Find User
	userModel, err := s.userRepo.GetUser(ctx, "id = ?", userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, err
	}

	// 2. Check if user is friend
	isFriend, err := s.friendRepo.CheckFriendExist(ctx, &model.Friend{
		UserId:   authenticatedUserId,
		FriendId: userModel.ID,
	})
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, err
	}

	if isFriend {
		userDto = mapper.MapUserToUserDtoWithoutSetting(userModel, consts.IS_FRIEND)

		return userDto, response.ErrCodeSuccess, http.StatusOK, nil
	}

	// 3. Check if user are send add friend request
	isSendFriendRequest, err := s.friendRequestRepo.CheckFriendRequestExist(ctx, &model.FriendRequest{
		UserId:   authenticatedUserId,
		FriendId: userModel.ID,
	})
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, err
	}

	if isSendFriendRequest {
		userDto = mapper.MapUserToUserDtoWithoutSetting(userModel, consts.SEND_FRIEND_REQUEST)

		return userDto, response.ErrCodeSuccess, http.StatusOK, nil
	}

	// 4. Check if user are receive add friend request
	isReceiveFriendRequest, _ := s.friendRequestRepo.CheckFriendRequestExist(ctx, &model.FriendRequest{
		UserId:   userModel.ID,
		FriendId: authenticatedUserId,
	})
	if isReceiveFriendRequest {
		userDto = mapper.MapUserToUserDtoWithoutSetting(userModel, consts.RECEIVE_FRIEND_REQUEST)

		return userDto, response.ErrCodeSuccess, http.StatusOK, nil
	}

	userDto = mapper.MapUserToUserDtoWithoutSetting(userModel, consts.NOT_FRIEND)

	return userDto, response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sUserInfo) GetManyUsers(
	ctx context.Context,
	query *query_object.UserQueryObject,
) (users []*model.User, resultCode int, httpStatusCode int, pagingResponse *response.PagingResponse, err error) {
	userModels, paging, err := s.userRepo.GetManyUser(ctx, query)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, nil, err
		}
		return nil, response.ErrDataNotFound, http.StatusInternalServerError, nil, err
	}

	return userModels, response.ErrCodeSuccess, http.StatusOK, paging, nil
}

func (s *sUserInfo) UpdateUser(
	ctx context.Context,
	userId uuid.UUID,
	updateData map[string]interface{},
	inAvatarUrl multipart.File,
	inCapwallUrl multipart.File,
	language consts.Language,
) (user *model.User, resultCode int, httpStatusCode int, err error) {
	// 1. update setting language
	if language != "" {
		settingFound, err := s.settingRepo.GetSetting(ctx, "user_id=?", userId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, response.ErrDataNotFound, http.StatusBadRequest, err
			}
			return nil, response.ErrDataNotFound, http.StatusInternalServerError, fmt.Errorf("Failed to get setting for user %v: %w", userId, err)
		}
		_, err = s.settingRepo.UpdateSetting(ctx, settingFound.ID, map[string]interface{}{
			"language": language,
		})
	}

	// 2. update user information
	userModel, err := s.userRepo.UpdateUser(ctx, userId, updateData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, response.ErrDataNotFound, http.StatusInternalServerError, err
	}

	// 3. update Avatar
	if inAvatarUrl != nil {
		avatarUrl, err := cloudinary_util.UploadMediaToCloudinary(inAvatarUrl)
		if err != nil {
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to upload Avatar: %w", err)
		}

		userModel.AvatarUrl = avatarUrl

		_, err = s.userRepo.UpdateUser(ctx, userId, map[string]interface{}{
			"avatar_url": userModel.AvatarUrl,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, response.ErrDataNotFound, http.StatusBadRequest, err
			}
			return nil, response.ErrServerFailed, http.StatusInternalServerError, err
		}
	}

	// 4. update Capwall
	if inCapwallUrl != nil {
		capwallUrl, err := cloudinary_util.UploadMediaToCloudinary(inCapwallUrl)
		if err != nil {
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to upload Capwall: %w", err)
		}

		userModel.CapwallUrl = capwallUrl

		_, err = s.userRepo.UpdateUser(ctx, userId, map[string]interface{}{
			"capwall_url": userModel.CapwallUrl,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, response.ErrDataNotFound, http.StatusBadRequest, err
			}
			return nil, response.ErrServerFailed, http.StatusInternalServerError, err
		}
	}

	return userModel, response.ErrCodeSuccess, http.StatusOK, nil
}
