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

func FromLikeUserPostModel(
	likeUserPost *models.LikeUserPost,
) *entities.LikeUserPost {
	if likeUserPost == nil {
		return nil
	}

	return &entities.LikeUserPost{
		UserId: likeUserPost.UserId,
		PostId: likeUserPost.PostId,
	}
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
