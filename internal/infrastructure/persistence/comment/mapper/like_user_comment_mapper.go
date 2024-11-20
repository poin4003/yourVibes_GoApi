package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToLikeUserCommentModel(
	likeUserComment *entities.LikeUserComment,
) *models.LikeUserComment {
	luc := &models.LikeUserComment{
		UserId:    likeUserComment.UserId,
		CommentId: likeUserComment.CommentId,
	}

	return luc
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
