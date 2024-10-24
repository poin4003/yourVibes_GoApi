package service_implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
) (resultCode int, httpStatusCode int, err error) {
	_, err = s.postRepo.GetPost(ctx, "id=?", likeUserPostModel.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find post %w", err.Error())
	}

	checkLike, err := s.postLikeRepo.CheckUserLikePost(ctx, likeUserPostModel)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to check like: %w", err)
	}
	if !checkLike {
		if err := s.postLikeRepo.CreateLikeUserPost(ctx, likeUserPostModel); err != nil {
			return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create like: %w", err)
		}
		return response.ErrCodeSuccess, http.StatusOK, nil
	} else {
		if err := s.postLikeRepo.DeleteLikeUserPost(ctx, likeUserPostModel); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ErrDataNotFound, http.StatusBadRequest, fmt.Errorf("failed to find delete like: %w", err)
			}
			return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete like: %w", err)
		}
		return response.ErrCodeSuccess, http.StatusNoContent, nil
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
