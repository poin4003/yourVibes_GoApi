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

type sCommentLike struct {
	userRepo            repository.IUserRepository
	commentRepo         repository.ICommentRepository
	likeUserCommentRepo repository.ILikeUserCommentRepository
}

func NewCommentLikeImplement(
	userRepo repository.IUserRepository,
	commentRepo repository.ICommentRepository,
	likeUserCommentRepo repository.ILikeUserCommentRepository,
) *sCommentLike {
	return &sCommentLike{
		userRepo:            userRepo,
		commentRepo:         commentRepo,
		likeUserCommentRepo: likeUserCommentRepo,
	}
}

func (s *sCommentLike) LikeComment(
	ctx context.Context,
	likeUserComment *model.LikeUserComment,
) (comment *model.Comment, resultCode int, httpStatusCode int, err error) {
	commentFound, err := s.commentRepo.GetOneComment(ctx, "id=?", likeUserComment.CommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find comment %w", err.Error())
	}

	checkLikeComment, err := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserComment)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to check like: %w", err)
	}

	if !checkLikeComment {
		if err := s.likeUserCommentRepo.CreateLikeUserComment(ctx, likeUserComment); err != nil {
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create like: %w", err)
		}

		commentFound.LikeCount++

		_, err = s.commentRepo.UpdateOneComment(ctx, commentFound.ID, map[string]interface{}{
			"like_count": commentFound.LikeCount,
		})

		return commentFound, response.ErrCodeSuccess, http.StatusOK, nil
	} else {
		if err := s.likeUserCommentRepo.DeleteLikeUserComment(ctx, likeUserComment); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, response.ErrDataNotFound, http.StatusBadRequest, fmt.Errorf("failed to find delete like: %w", err)
			}
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete like: %w", err)
		}

		commentFound.LikeCount--

		_, err = s.commentRepo.UpdateOneComment(ctx, commentFound.ID, map[string]interface{}{
			"like_count": commentFound.LikeCount,
		})

		return commentFound, response.ErrCodeSuccess, http.StatusNoContent, nil
	}
}

func (s *sCommentLike) GetUsersOnLikeComment(
	ctx context.Context,
	commentId uuid.UUID,
	query *query_object.CommentLikeQueryObject,
) (users []*model.User, resultCode int, httpStatusCode int, responsePaging *response.PagingResponse, err error) {
	likeUserComment, paging, err := s.likeUserCommentRepo.GetLikeUserComment(ctx, commentId, query)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, nil, err
	}

	return likeUserComment, response.ErrCodeSuccess, http.StatusOK, paging, nil
}
