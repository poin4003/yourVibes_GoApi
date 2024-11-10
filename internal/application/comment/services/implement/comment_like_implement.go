package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	comment_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
)

type sCommentLike struct {
	userRepo            comment_repo.IUserRepository
	commentRepo         comment_repo.ICommentRepository
	likeUserCommentRepo comment_repo.ILikeUserCommentRepository
}

func NewCommentLikeImplement(
	userRepo comment_repo.IUserRepository,
	commentRepo comment_repo.ICommentRepository,
	likeUserCommentRepo comment_repo.ILikeUserCommentRepository,
) *sCommentLike {
	return &sCommentLike{
		userRepo:            userRepo,
		commentRepo:         commentRepo,
		likeUserCommentRepo: likeUserCommentRepo,
	}
}

func (s *sCommentLike) LikeComment(
	ctx context.Context,
	likeUserComment *models.LikeUserComment,
	userId uuid.UUID,
) (commentDto *response.CommentDto, resultCode int, httpStatusCode int, err error) {
	commentFound, err := s.commentRepo.GetOneComment(ctx, "id=?", likeUserComment.CommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg_response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find comment %w", err.Error())
	}

	checkLikeComment, err := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserComment)
	if err != nil {
		return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to check like: %w", err)
	}

	if !checkLikeComment {
		if err := s.likeUserCommentRepo.CreateLikeUserComment(ctx, likeUserComment); err != nil {
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create like: %w", err)
		}

		commentFound.LikeCount++

		_, err = s.commentRepo.UpdateOneComment(ctx, commentFound.ID, map[string]interface{}{
			"like_count": commentFound.LikeCount,
		})

		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, &models.LikeUserComment{
			CommentId: commentFound.ID,
			UserId:    userId,
		})

		commentDto = mapper.MapCommentToCommentDto(commentFound, isLiked)

		return commentDto, pkg_response.ErrCodeSuccess, http.StatusOK, nil
	} else {
		if err := s.likeUserCommentRepo.DeleteLikeUserComment(ctx, likeUserComment); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg_response.ErrDataNotFound, http.StatusBadRequest, fmt.Errorf("failed to find delete like: %w", err)
			}
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete like: %w", err)
		}

		commentFound.LikeCount--

		_, err = s.commentRepo.UpdateOneComment(ctx, commentFound.ID, map[string]interface{}{
			"like_count": commentFound.LikeCount,
		})

		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, &models.LikeUserComment{
			CommentId: commentFound.ID,
			UserId:    userId,
		})

		commentDto = mapper.MapCommentToCommentDto(commentFound, isLiked)

		return commentDto, pkg_response.ErrCodeSuccess, http.StatusNoContent, nil
	}
}

func (s *sCommentLike) GetUsersOnLikeComment(
	ctx context.Context,
	commentId uuid.UUID,
	query *query.CommentLikeQueryObject,
) (users []*models.User, resultCode int, httpStatusCode int, pkg_responsePaging *pkg_response.PagingResponse, err error) {
	likeUserComment, paging, err := s.likeUserCommentRepo.GetLikeUserComment(ctx, commentId, query)
	if err != nil {
		return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, nil, err
	}

	return likeUserComment, pkg_response.ErrCodeSuccess, http.StatusOK, paging, nil
}
