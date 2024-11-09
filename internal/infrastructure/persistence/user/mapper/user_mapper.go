package mapper

import (
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToUserModel(user *user_entity.User) *models.User {
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

func FromUserModel(userModel *models.User) *user_entity.User {
	var setting = &user_entity.Setting{
		ID:        userModel.Setting.ID,
		UserId:    userModel.Setting.UserId,
		Language:  userModel.Setting.Language,
		Status:    userModel.Setting.Status,
		CreatedAt: userModel.Setting.CreatedAt,
		UpdatedAt: userModel.Setting.UpdatedAt,
	}

	var user = &user_entity.User{
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
	user.ID = user.ID

	return user
}

func FromUserModelList(userModelList []*models.User) []*user_entity.User {
	userEntityList := make([]*user_entity.User, len(userModelList))
	for i, userModel := range userModelList {
		userEntityList[i] = FromUserModel(userModel)
	}

	return userEntityList
}
