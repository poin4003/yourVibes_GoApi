package implement

import (
	"context"

	commentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	commentEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	commentRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
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
	result = &commentCommand.LikeCommentResult{
		Comment: nil,
	}
	// 1. Get comment by id
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if commentFound == nil {
		return nil, response.NewDataNotFoundError("comment not found")
	}

	// 2. Check status of like
	likeUserCommentEntity, err := commentEntity.NewLikeUserCommentEntity(command.UserId, command.CommentId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	checkLikeComment, err := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if !checkLikeComment {
		// 2.1. Create like if not exits
		if err := s.likeUserCommentRepo.CreateLikeUserComment(ctx, likeUserCommentEntity); err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 2.2. Plus 1 to like count of comment
		updateCommentData := commentEntity.CommentUpdate{
			LikeCount: pointer.Ptr(commentFound.LikeCount + 1),
		}

		err = updateCommentData.ValidateCommentUpdate()
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		s.commentRepo.UpdateOne(ctx, commentFound.ID, &updateCommentData)

		// 2. Check like status of authenticated user
		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)

		result.Comment = mapper.NewCommentWithLikedResultFromEntityAndIsLiked(commentFound, isLiked)

		return result, nil
	} else {
		// 3.1. Delete like if it exits
		if err := s.likeUserCommentRepo.DeleteLikeUserComment(ctx, likeUserCommentEntity); err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 3.2. Minus 1 of comment like count
		updateCommentData := commentEntity.CommentUpdate{
			LikeCount: pointer.Ptr(commentFound.LikeCount - 1),
		}

		err = updateCommentData.ValidateCommentUpdate()
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		s.commentRepo.UpdateOne(ctx, commentFound.ID, &updateCommentData)

		// 3.3. Check like status of authenticated user
		isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)

		result.Comment = mapper.NewCommentWithLikedResultFromEntityAndIsLiked(commentFound, isLiked)

		return result, nil
	}
}

func (s *sCommentLike) GetUsersOnLikeComment(
	ctx context.Context,
	query *commentQuery.GetCommentLikeQuery,
) (result *commentQuery.GetCommentLikeResult, err error) {
	likeUserCommentEntites, paging, err := s.likeUserCommentRepo.GetLikeUserComment(ctx, query)
	if err != nil {
		return result, err
	}

	var likeUserCommentResults []*common.UserResult
	for _, likeUserCommentEntity := range likeUserCommentEntites {
		likeUserCommentResults = append(likeUserCommentResults, mapper.NewUserResultFromEntity(likeUserCommentEntity))
	}

	return &commentQuery.GetCommentLikeResult{
		Users:          likeUserCommentResults,
		PagingResponse: paging,
	}, nil
}
