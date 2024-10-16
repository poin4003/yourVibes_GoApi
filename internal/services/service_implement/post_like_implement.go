package service_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
)

type sPostLike struct {
	userRepo repository.IUserRepository
	postRepo repository.IPostRepository
}

func NewPostLikeImplement(
	userRepo repository.IUserRepository,
	postRepo repository.IPostRepository,
) *sPostLike {
	return &sPostLike{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (s *sPostLike) LikePost(
	ctx context.Context,
	likeUserPost *model.LikeUserPost,
) error {
	return nil
}

func (s *sPostLike) DeleteLikePost(
	ctx context.Context,
	likeUserPost *model.LikeUserPost,
) error {
	return nil
}

func (s *sPostLike) GetUsersOnLikes(
	ctx context.Context,
	postId uuid.UUID,
) ([]*model.User, error) {
	return []*model.User{}, nil
}
