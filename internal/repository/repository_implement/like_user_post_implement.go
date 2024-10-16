package repository_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"gorm.io/gorm"
)

type rLikeUserPost struct {
	db *gorm.DB
}

func NewLikeUserPostRepositoryImplement(db *gorm.DB) *rLikeUserPost {
	return &rLikeUserPost{db: db}
}

func (r *rLikeUserPost) CreateLikeUserPost(
	ctx context.Context,
	likeUserPost *model.LikeUserPost,
) error {
	return nil
}

func (r *rLikeUserPost) DeleteLikeUserPost(
	ctx context.Context,
	likeUserPost *model.LikeUserPost,
) error {
	return nil
}

func (r *rLikeUserPost) GetLikeUserPost(
	ctx context.Context,
	postId uuid.UUID,
) ([]*model.User, error) {
	return []*model.User{}, nil
}
