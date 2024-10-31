package service_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/post_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type sUserNewFeed struct {
	userRepo         repository.IUserRepository
	postRepo         repository.IPostRepository
	likeUserPostRepo repository.ILikeUserPostRepository
	newFeedRepo      repository.INewFeedRepository
}

func NewUserNewFeedImplement(
	userRepo repository.IUserRepository,
	postRepo repository.IPostRepository,
	likeUserPostRepo repository.ILikeUserPostRepository,
	newFeedRepo repository.INewFeedRepository,
) *sUserNewFeed {
	return &sUserNewFeed{
		userRepo:         userRepo,
		postRepo:         postRepo,
		likeUserPostRepo: likeUserPostRepo,
		newFeedRepo:      newFeedRepo,
	}
}

func (s *sUserNewFeed) DeleteNewFeed(
	ctx context.Context,
	userId uuid.UUID,
	postId uuid.UUID,
) (resultCode int, httpStatusCode int, err error) {
	err = s.newFeedRepo.DeleteNewFeed(ctx, userId, postId)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, err
	}

	return response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sUserNewFeed) GetNewFeeds(
	ctx context.Context,
	userId uuid.UUID,
	query *query_object.NewFeedQueryObject,
) (postDtos []*post_dto.PostDto, pagingResponse *response.PagingResponse, resultCode int, httpStatusCode int, err error) {
	postModels, paging, err := s.newFeedRepo.GetManyNewFeed(ctx, userId, query)
	if err != nil {
		return nil, nil, response.ErrServerFailed, http.StatusInternalServerError, err
	}

	for _, post := range postModels {
		isLiked, _ := s.likeUserPostRepo.CheckUserLikePost(ctx, &model.LikeUserPost{
			PostId: post.ID,
			UserId: userId,
		})

		postDto := mapper.MapPostToPostDto(post, isLiked)
		postDtos = append(postDtos, postDto)
	}

	return postDtos, paging, response.ErrCodeSuccess, http.StatusOK, nil
}
