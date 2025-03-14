package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToLikeUserPostModel(
	likeUserPost *entities.LikeUserPost,
) *models.LikeUserPost {
	lup := &models.LikeUserPost{
		UserId: likeUserPost.UserId,
		PostId: likeUserPost.PostId,
	}

	return lup
}

func FromUserModel(userModel *models.User) *entities.User {
	var user = &entities.User{
		FamilyName: userModel.FamilyName,
		Name:       userModel.Name,
		AvatarUrl:  userModel.AvatarUrl,
	}

	user.ID = userModel.ID

	return user
}
