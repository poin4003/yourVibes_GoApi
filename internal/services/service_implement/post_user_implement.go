package service_implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/internal/vo"
)

type sPostUser struct {
	repo repository.IPostRepository
}

func NewPostUserImplement(repo repository.IPostRepository) *sPostUser {
	return &sPostUser{repo: repo}
}

func (s *sPostUser) CreatePost(
	ctx context.Context,
	in *vo.CreatePostInput,
) (post *model.Post, resultCode int, err error) {
	return &model.Post{}, 0, nil
}

func (s *sPostUser) UpdatePost(
	ctx context.Context,
	in *vo.UpdatePostInput,
) (post *model.Post, resultCode int, err error) {
	return &model.Post{}, 0, nil
}

func (s *sPostUser) DeletePost(
	ctx context.Context,
	email string) (resultCode int, err error) {
	return 0, nil
}

func (s *sPostUser) GetPost(
	ctx context.Context,
	query interface{}, args ...interface{},
) (post *model.Post, resultCode int, err error) {
	return &model.Post{}, 0, nil
}

func (s *sPostUser) GetAllPost(
	ctx context.Context,
) (posts []*model.Post, resultCode int, err error) {
	return []*model.Post{}, 0, nil
}
