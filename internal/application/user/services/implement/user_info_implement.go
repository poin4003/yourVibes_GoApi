package implement

import (
	"context"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/media"

	"github.com/google/uuid"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	userMapper "github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	userQuery "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	userRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
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
		User:       nil,
		ResultCode: response2.ErrServerFailed,
	}
	// 1. Find User
	userFound, err := s.userRepo.GetOne(ctx, "id = ?", query.UserId)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return nil, response2.NewDataNotFoundError("user not found")
	}

	// 2. Return if user fetches his own information
	if query.AuthenticatedUserId == query.UserId {
		result.User = userMapper.NewUserResultWithoutSettingEntity(userFound, consts.NOT_FRIEND)
		return result, nil
	}

	// 3. Check friend status
	var friendStatus consts.FriendStatus
	isFriend, err := s.friendRepo.CheckFriendExist(ctx, &userEntity.Friend{
		UserId:   query.AuthenticatedUserId,
		FriendId: query.UserId,
	})
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
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
			return nil, response2.NewServerFailedError(err.Error())
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
				return nil, response2.NewServerFailedError(err.Error())
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
		resultCode = response2.ErrCodeSuccess
	case consts.FRIEND_ONLY:
		if friendStatus == consts.IS_FRIEND {
			userResult = userMapper.NewUserResultWithoutSettingEntity(userFound, friendStatus)
			resultCode = response2.ErrCodeSuccess
		} else {
			userResult = userMapper.NewUserResultWithoutPrivateInfo(userFound, friendStatus)
			resultCode = response2.ErrUserFriendAccess
		}
	case consts.PRIVATE:
		userResult = userMapper.NewUserResultWithoutPrivateInfo(userFound, friendStatus)
		resultCode = response2.ErrUserPrivateAccess
	default:
		userResult = userMapper.NewUserResultWithoutPrivateInfo(userFound, friendStatus)
		resultCode = response2.ErrUserPrivateAccess
	}

	result.User = userResult
	result.ResultCode = resultCode
	return result, nil
}

func (s *sUserInfo) GetManyUsers(
	ctx context.Context,
	query *userQuery.GetManyUserQuery,
) (result *userQuery.UserQueryListResult, err error) {
	userEntities, paging, err := s.userRepo.GetMany(ctx, query)

	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	var userResultList []*common.UserShortVerResult
	for _, user := range userEntities {
		userResultList = append(userResultList, userMapper.NewUserShortVerEntity(user))
	}

	return &userQuery.UserQueryListResult{
		Users:          userResultList,
		PagingResponse: paging,
	}, nil
}

func (s *sUserInfo) UpdateUser(
	ctx context.Context,
	command *userCommand.UpdateUserCommand,
) (result *userCommand.UpdateUserCommandResult, err error) {
	// 1. find user
	userFound, err := s.userRepo.GetById(ctx, *command.UserId)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return nil, response2.NewDataNotFoundError("user not found")
	}

	// 1. update setting language
	if command.LanguageSetting != nil {
		settingFound, err := s.settingRepo.GetSetting(ctx, "user_id=?", command.UserId)
		if err != nil {
			return nil, response2.NewServerFailedError(err.Error())
		}

		if settingFound == nil {
			return nil, response2.NewDataNotFoundError("setting not found")
		}

		s.settingRepo.UpdateOne(ctx, settingFound.ID,
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
		return nil, response2.NewServerFailedError(err.Error())
	}

	// 3. update Avatar
	if command.Avatar != nil && command.Avatar.Size > 0 && command.Avatar.Filename != "" {
		avatarUrl, err := media.SaveMedia(command.Avatar)
		if err != nil {
			return nil, response2.NewServerFailedError(err.Error())
		}

		_, err = s.userRepo.UpdateOne(ctx, *command.UserId, &userEntity.UserUpdate{
			AvatarUrl: &avatarUrl,
		})
		if err != nil {
			return nil, response2.NewServerFailedError(err.Error())
		}
	}

	// 4. update Capwall
	if command.Capwall != nil && command.Capwall.Size > 0 && command.Capwall.Filename != "" {
		capwallUrl, err := media.SaveMedia(command.Capwall)
		if err != nil {
			return nil, response2.NewServerFailedError(err.Error())
		}

		_, err = s.userRepo.UpdateOne(ctx, *command.UserId, &userEntity.UserUpdate{
			CapwallUrl: &capwallUrl,
		})
		if err != nil {
			return nil, response2.NewServerFailedError(err.Error())
		}
	}

	userFound, err = s.userRepo.UpdateOne(ctx, *command.UserId, updateUserEntity)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	return &userCommand.UpdateUserCommandResult{
		User: userMapper.NewUserResultFromEntity(userFound),
	}, nil
}

func (s *sUserInfo) GetUserStatusById(
	ctx context.Context,
	id uuid.UUID,
) (status bool, err error) {
	userStatus, err := s.userRepo.GetStatusById(ctx, id)
	if err != nil {
		return false, response2.NewServerFailedError(err.Error())
	}
	return userStatus, nil
}
