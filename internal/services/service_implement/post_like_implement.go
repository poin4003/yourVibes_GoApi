package service_implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/post_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
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
	userId uuid.UUID,
) (postDto *post_dto.PostDto, resultCode int, httpStatusCode int, err error) {
	postFound, err := s.postRepo.GetPost(ctx, "id=?", likeUserPostModel.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find post %w", err.Error())
	}

	checkLiked, err := s.postLikeRepo.CheckUserLikePost(ctx, likeUserPostModel)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to check like: %w", err)
	}

	if !checkLiked {
		if err := s.postLikeRepo.CreateLikeUserPost(ctx, likeUserPostModel); err != nil {
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create like: %w", err)
		}

		postFound.LikeCount++

		_, err = s.postRepo.UpdatePost(ctx, postFound.ID, map[string]interface{}{
			"like_count": postFound.LikeCount,
		})

		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, &model.LikeUserPost{
			PostId: postFound.ID,
			UserId: userId,
		})

		postDto = mapper.MapPostToPostDto(postFound, isLiked)

		return postDto, response.ErrCodeSuccess, http.StatusOK, nil
	} else {
		if err := s.postLikeRepo.DeleteLikeUserPost(ctx, likeUserPostModel); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, response.ErrDataNotFound, http.StatusBadRequest, fmt.Errorf("failed to find delete like: %w", err)
			}
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete like: %w", err)
		}

		postFound.LikeCount--

		_, err = s.postRepo.UpdatePost(ctx, postFound.ID, map[string]interface{}{
			"like_count": postFound.LikeCount,
		})

		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, &model.LikeUserPost{
			PostId: postFound.ID,
			UserId: userId,
		})

		postDto = mapper.MapPostToPostDto(postFound, isLiked)

		return postDto, response.ErrCodeSuccess, http.StatusOK, nil
	}
}

func (s *sPostLike) GetUsersOnLikes(
	ctx context.Context,
	postId uuid.UUID,
	query *query_object.PostLikeQueryObject,
) (users []*model.User, resultCode int, httpStatusCode int, responsePaging *response.PagingResponse, err error) {
	likeUserPostModel, paging, err := s.postLikeRepo.GetLikeUserPost(ctx, postId, query)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, nil, err
	}

	return likeUserPostModel, response.ErrCodeSuccess, http.StatusOK, paging, nil
}
