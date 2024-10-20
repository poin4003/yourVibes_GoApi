package service_implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type sPostLike struct {
	userRepo     repository.IUserRepository
	postRepo     repository.IPostRepository
	postLikeRepo repository.ILikeUserPostRepository
}

func NewPostLikeImplement(
	userRepo repository.IUserRepository,
	postRepo repository.IPostRepository,
	postLikeRepo repository.ILikeUserPostRepository,
) *sPostLike {
	return &sPostLike{
		userRepo:     userRepo,
		postRepo:     postRepo,
		postLikeRepo: postLikeRepo,
	}
}

func (s *sPostLike) LikePost(
	ctx context.Context,
	likeUserPostModel *model.LikeUserPost,
) (resultCode int, err error) {
	if err := s.postLikeRepo.CreateLikeUserPost(ctx, likeUserPostModel); err != nil {
		return response.ErrServerFailed, err
	}
	return response.ErrCodeSuccess, nil
}

func (s *sPostLike) DeleteLikePost(
	ctx context.Context,
	likeUserPostModel *model.LikeUserPost,
) (resultCode int, err error) {
	if err := s.postLikeRepo.DeleteLikeUserPost(ctx, likeUserPostModel); err != nil {
		return response.ErrServerFailed, err
	}
	return response.ErrCodeSuccess, nil

}

func (s *sPostLike) GetUsersOnLikes(
	ctx context.Context,
	query *query_object.PostLikeQueryObject,
) (users []*model.User, resultCode int, err error) {
	likeUserPostModel, err := s.postLikeRepo.GetLikeUserPost(ctx, query)
	if err != nil {
		return nil, response.ErrDataNotFound, err
	}

	return likeUserPostModel, response.ErrCodeSuccess, nil
}
