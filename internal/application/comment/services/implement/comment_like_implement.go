package implement

import (
	"context"
	"errors"
	"fmt"
	commentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	commentEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	commentRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
	"net/http"
)

type sCommentLike struct {
	userRepo            commentRepo.IUserRepository
	commentRepo         commentRepo.ICommentRepository
	likeUserCommentRepo commentRepo.ILikeUserCommentRepository
}

func NewCommentLikeImplement(
	userRepo commentRepo.IUserRepository,
	commentRepo commentRepo.ICommentRepository,
	likeUserCommentRepo commentRepo.ILikeUserCommentRepository,
) *sCommentLike {
	return &sCommentLike{
		userRepo:            userRepo,
		commentRepo:         commentRepo,
		likeUserCommentRepo: likeUserCommentRepo,
	}
}

func (s *sCommentLike) LikeComment(
	ctx context.Context,
	command *commentCommand.LikeCommentCommand,
) (result *commentCommand.LikeCommentResult, err error) {
	result = &commentCommand.LikeCommentResult{}
	result.Comment = nil
	result.ResultCode = response.ErrDataNotFound
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Get comment by id
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("error when find comment %v", err.Error())
	}

	// 2. Check status of like
	likeUserCommentEntity, err := commentEntity.NewLikeUserCommentEntity(command.UserId, command.CommentId)
	if err != nil {
		return result, err
	}

	checkLikeComment, err := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)
	if err != nil {
		return result, fmt.Errorf("failed to check like: %w", err)
	}

	if !checkLikeComment {
		// 2.1. Create like if not exits
		if err := s.likeUserCommentRepo.CreateLikeUserComment(ctx, likeUserCommentEntity); err != nil {
			return result, fmt.Errorf("failed to create like: %w", err)
		}

		// 2.2. Plus 1 to like count of comment
		updateCommentData := commentEntity.CommentUpdate{
			LikeCount: pointer.Ptr(commentFound.LikeCount + 1),
		}

		err = updateCommentData.ValidateCommentUpdate()
		if err != nil {
			return result, fmt.Errorf("failed to update comment: %w", err)
		}

		_, err = s.commentRepo.UpdateOne(ctx, commentFound.ID, &updateCommentData)

		// 2. Check like status of authenticated user
		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)

		result.Comment = mapper.NewCommentWithLikedResultFromEntityAndIsLiked(commentFound, isLiked)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	} else {
		// 3.1. Delete like if it exits
		if err := s.likeUserCommentRepo.DeleteLikeUserComment(ctx, likeUserCommentEntity); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, fmt.Errorf("failed to find delete like: %w", err)
			}
			return result, fmt.Errorf("failed to delete like: %w", err)
		}

		// 3.2. Minus 1 of comment like count
		updateCommentData := commentEntity.CommentUpdate{
			LikeCount: pointer.Ptr(commentFound.LikeCount - 1),
		}

		err = updateCommentData.ValidateCommentUpdate()
		if err != nil {
			return result, fmt.Errorf("failed to update comment: %w", err)
		}

		_, err = s.commentRepo.UpdateOne(ctx, commentFound.ID, &updateCommentData)

		// 3.3. Check like status of authenticated user
		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)

		result.Comment = mapper.NewCommentWithLikedResultFromEntityAndIsLiked(commentFound, isLiked)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}
}

func (s *sCommentLike) GetUsersOnLikeComment(
	ctx context.Context,
	query *commentQuery.GetCommentLikeQuery,
) (result *commentQuery.GetCommentLikeResult, err error) {
	result = &commentQuery.GetCommentLikeResult{}
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
