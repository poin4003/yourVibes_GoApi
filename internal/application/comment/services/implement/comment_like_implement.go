package implement

import (
	"context"
	"errors"
	"fmt"
	comment_command "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	comment_query "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	comment_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	comment_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
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
	command *comment_command.LikeCommentCommand,
) (result *comment_command.LikeCommentResult, err error) {
	result = &comment_command.LikeCommentResult{}
	// 1. Get comment by id
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Comment = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.Comment = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when find comment %w", err.Error())
	}

	// 2. Check status of like
	likeUserCommentEntity, err := comment_entity.NewLikeUserCommentEntity(command.UserId, command.CommentId)
	if err != nil {
		result.Comment = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	checkLikeComment, err := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)
	if err != nil {
		result.Comment = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to check like: %w", err)
	}

	if !checkLikeComment {
		// 2.1. Create like if not exits
		if err := s.likeUserCommentRepo.CreateLikeUserComment(ctx, likeUserCommentEntity); err != nil {
			result.Comment = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to create like: %w", err)
		}

		// 2.2. Plus 1 to like count of comment
		updateCommentData := comment_entity.CommentUpdate{
			LikeCount: pointer.Ptr(commentFound.LikeCount + 1),
		}

		err = updateCommentData.ValidateCommentUpdate()
		if err != nil {
			result.Comment = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to update comment: %w", err)
		}

		_, err = s.commentRepo.UpdateOne(ctx, commentFound.ID, &updateCommentData)

		// 2. Check like status of authenticated user
		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)

		result.Comment = mapper.NewCommentWithLikedResultFromEntity(commentFound, isLiked)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	} else {
		// 3.1. Delete like if it exits
		if err := s.likeUserCommentRepo.DeleteLikeUserComment(ctx, likeUserCommentEntity); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Comment = nil
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, fmt.Errorf("failed to find delete like: %w", err)
			}
			result.Comment = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to delete like: %w", err)
		}

		// 3.2. Minus 1 of comment like count
		updateCommentData := comment_entity.CommentUpdate{
			LikeCount: pointer.Ptr(commentFound.LikeCount - 1),
		}

		err = updateCommentData.ValidateCommentUpdate()
		if err != nil {
			result.Comment = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to update comment: %w", err)
		}

		_, err = s.commentRepo.UpdateOne(ctx, commentFound.ID, &updateCommentData)

		// 3.3. Check like status of authenticated user
		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)

		result.Comment = mapper.NewCommentWithLikedResultFromEntity(commentFound, isLiked)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}
}

func (s *sCommentLike) GetUsersOnLikeComment(
	ctx context.Context,
	query *comment_query.GetCommentLikeQuery,
) (result *comment_query.GetCommentLikeResult, err error) {
	result = &comment_query.GetCommentLikeResult{}
	likeUserCommentEntites, paging, err := s.likeUserCommentRepo.GetLikeUserComment(ctx, query)
	if err != nil {
		result.Users = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.PagingResponse = nil
		return result, err
	}

	var likeUserCommentResults []*common.UserResult
	for _, likeUserCommentEntity := range likeUserCommentEntites {
		likeUserCommentResults = append(likeUserCommentResults, mapper.NewUserResultFromEntity(likeUserCommentEntity))
	}

	result.Users = likeUserCommentResults
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}
