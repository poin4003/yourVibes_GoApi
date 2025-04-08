package mapper

import (
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToUserModel(user *userEntity.User) *models.User {
	u := &models.User{
		FamilyName:   user.FamilyName,
		Name:         user.Name,
		Email:        user.Email,
		Password:     user.Password,
		PhoneNumber:  user.PhoneNumber,
		Birthday:     user.Birthday,
		AvatarUrl:    user.AvatarUrl,
		CapwallUrl:   user.CapwallUrl,
		Privacy:      user.Privacy,
		Biography:    user.Biography,
		AuthType:     user.AuthType,
		AuthGoogleId: user.AuthGoogleId,
		PostCount:    user.PostCount,
		FriendCount:  user.FriendCount,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
	u.ID = user.ID

	return u
}

func FromUserModel(userModel *models.User) *userEntity.User {
	var setting = &userEntity.Setting{
		ID:        userModel.Setting.ID,
		UserId:    userModel.Setting.UserId,
		Language:  userModel.Setting.Language,
		Status:    userModel.Setting.Status,
		CreatedAt: userModel.Setting.CreatedAt,
		UpdatedAt: userModel.Setting.UpdatedAt,
	}

	var user = &userEntity.User{
		FamilyName:   userModel.FamilyName,
		Name:         userModel.Name,
		Email:        userModel.Email,
		Password:     userModel.Password,
		PhoneNumber:  userModel.PhoneNumber,
		Birthday:     userModel.Birthday,
		AvatarUrl:    userModel.AvatarUrl,
		CapwallUrl:   userModel.CapwallUrl,
		Privacy:      userModel.Privacy,
		Biography:    userModel.Biography,
		AuthType:     userModel.AuthType,
		AuthGoogleId: userModel.AuthGoogleId,
		PostCount:    userModel.PostCount,
		FriendCount:  userModel.FriendCount,
		Status:       userModel.Status,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
		Setting:      setting,
	}
	user.ID = userModel.ID

	return user
}

func FromUserModelList(userModelList []*models.User) []*userEntity.User {
	var userEntityList []*userEntity.User
	for _, userModel := range userModelList {
		userEntityList = append(userEntityList, FromUserModel(userModel))
	}

	return userEntityList
}

func FromUserModelWithSendFriendRequest(
	userModel *models.User,
	isSendFriendRequest bool,
) *userEntity.UserWithSendFriendRequest {
	var user = &userEntity.UserWithSendFriendRequest{
		FamilyName:          userModel.FamilyName,
		Name:                userModel.Name,
		AvatarUrl:           userModel.AvatarUrl,
		IsSendFriendRequest: isSendFriendRequest,
	}
	user.ID = userModel.ID

	return user
}
