package service_implement

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/cloudinary_util"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"mime/multipart"
	"net/http"

	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
)

type sUserInfo struct {
	userRepo repository.IUserRepository
}

func NewUserInfoImplement(userRepo repository.IUserRepository) *sUserInfo {
	return &sUserInfo{userRepo: userRepo}
}

func (s *sUserInfo) GetInfoByUserId(
	ctx context.Context,
	userId uuid.UUID,
) (user *model.User, resultCode int, err error) {
	userModel, err := s.userRepo.GetUser(ctx, "id = ?", userId)

	if err != nil {
		return nil, response.ErrDataNotFound, err
	}

	return userModel, response.ErrCodeSuccess, nil
}

func (s *sUserInfo) GetManyUsers(
	ctx context.Context,
	query *query_object.UserQueryObject,
) (users []*model.User, resultCode int, err error) {
	userModels, err := s.userRepo.GetManyUser(ctx, query)

	if err != nil {
		return nil, response.ErrDataNotFound, err
	}

	return userModels, response.ErrCodeSuccess, nil
}

func (s *sUserInfo) UpdateUser(
	ctx context.Context,
	userId uuid.UUID,
	updateData map[string]interface{},
	inAvatarUrl multipart.File,
	inCapwallUrl multipart.File,
) (user *model.User, resultCode int, err error) {
	// 1. update user information
	userModel, err := s.userRepo.UpdateUser(ctx, userId, updateData)
	if err != nil {
		return nil, response.ErrDataNotFound, err
	}

	// 2. update Avatar
	if inAvatarUrl != nil {
		avatarUrl, err := cloudinary_util.UploadMediaToCloudinary(inAvatarUrl)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to upload Avatar: %w", err)
		}

		userModel.AvatarUrl = avatarUrl

		_, err = s.userRepo.UpdateUser(ctx, userId, map[string]interface{}{
			"avatar_url": userModel.AvatarUrl,
		})
	}

	// 3. update Capwall
	if inCapwallUrl != nil {
		capwallUrl, err := cloudinary_util.UploadMediaToCloudinary(inCapwallUrl)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to upload Capwall: %w", err)
		}

		userModel.CapwallUrl = capwallUrl

		_, err = s.userRepo.UpdateUser(ctx, userId, map[string]interface{}{
			"capwall_url": userModel.CapwallUrl,
		})
	}

	return userModel, response.ErrCodeSuccess, nil
}
