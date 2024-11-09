package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	comment_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repository"
	entities2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/comment/comment_user/dto/mapper"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/comment/comment_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/comment/comment_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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
	likeUserComment *entities2.LikeUserComment,
	userId uuid.UUID,
) (commentDto *response2.CommentDto, resultCode int, httpStatusCode int, err error) {
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

		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, &entities2.LikeUserComment{
			CommentId: commentFound.ID,
			UserId:    userId,
		})

		commentDto = mapper.MapCommentToCommentDto(commentFound, isLiked)

		return commentDto, response.ErrCodeSuccess, http.StatusOK, nil
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

		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, &entities2.LikeUserComment{
			CommentId: commentFound.ID,
			UserId:    userId,
		})

		commentDto = mapper.MapCommentToCommentDto(commentFound, isLiked)

		return commentDto, response.ErrCodeSuccess, http.StatusNoContent, nil
	}
}

func (s *sCommentLike) GetUsersOnLikeComment(
	ctx context.Context,
	commentId uuid.UUID,
	query *query.CommentLikeQueryObject,
) (users []*entities2.User, resultCode int, httpStatusCode int, responsePaging *response.PagingResponse, err error) {
	likeUserComment, paging, err := s.likeUserCommentRepo.GetLikeUserComment(ctx, commentId, query)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, nil, err
	}

	return likeUserComment, response.ErrCodeSuccess, http.StatusOK, paging, nil
}
