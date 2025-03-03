package implement

import (
	"context"

	"github.com/google/uuid"
	commentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	commentEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	commentValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/validator"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	commentRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
)

type sCommentUser struct {
	commentRepo         commentRepo.ICommentRepository
	userRepo            commentRepo.IUserRepository
	postRepo            commentRepo.IPostRepository
	likeUserCommentRepo commentRepo.ILikeUserCommentRepository
	commentReportRepor  commentRepo.ICommentReportRepository
}

func NewCommentUserImplement(
	commentRepo commentRepo.ICommentRepository,
	userRepo commentRepo.IUserRepository,
	postRepo commentRepo.IPostRepository,
	likeUserCommentRepo commentRepo.ILikeUserCommentRepository,
	commentReportRepo commentRepo.ICommentReportRepository,
) *sCommentUser {
	return &sCommentUser{
		commentRepo:         commentRepo,
		userRepo:            userRepo,
		postRepo:            postRepo,
		likeUserCommentRepo: likeUserCommentRepo,
		commentReportRepor:  commentReportRepo,
	}
}

func (s *sCommentUser) CreateComment(
	ctx context.Context,
	command *commentCommand.CreateCommentCommand,
) (result *commentCommand.CreateCommentResult, err error) {
	// 1. Find post
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if postFound == nil {
		return nil, response.NewDataNotFoundError("post not found")
	}

	if command.ParentId != nil {
		// 2.1. Get root comment
		parentComment, err := s.commentRepo.GetById(ctx, *command.ParentId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		if parentComment == nil {
			return nil, response.NewDataNotFoundError("parent comment not found")
		}

		// 2. Update rep count +1
		updateComment := &commentEntity.CommentUpdate{
			RepCommentCount: pointer.Ptr(parentComment.RepCommentCount + 1),
		}

		if err = updateComment.ValidateCommentUpdate(); err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		_, err = s.commentRepo.UpdateOne(ctx, parentComment.ID, updateComment)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}
	}

	// 4. Create a comment
	newComment, _ := commentEntity.NewComment(
		command.PostId,
		command.UserId,
		command.ParentId,
		command.Content,
	)

	commentCreated, err := s.commentRepo.CreateOne(ctx, newComment)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 5. Update comment count for post
	updatePost := &postEntity.PostUpdate{
		CommentCount: pointer.Ptr(postFound.CommentCount + 1),
	}

	err = updatePost.ValidatePostUpdate()
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePost)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 6. Validate comment after create
	validateComment, err := commentValidator.NewValidatedComment(commentCreated)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return &commentCommand.CreateCommentResult{
		Comment: mapper.NewCommentResultFromValidateEntity(validateComment),
	}, nil
}

func (s *sCommentUser) UpdateComment(
	ctx context.Context,
	command *commentCommand.UpdateCommentCommand,
) (result *commentCommand.UpdateCommentResult, err error) {
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if commentFound == nil {
		return nil, response.NewDataNotFoundError("comment not found")
	}

	updateData := &commentEntity.CommentUpdate{
		Content: command.Content,
	}

	err = updateData.ValidateCommentUpdate()
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	commentUpdate, err := s.commentRepo.UpdateOne(ctx, command.CommentId, updateData)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return &commentCommand.UpdateCommentResult{
		Comment: mapper.NewCommentResultFromEntity(commentUpdate),
	}, nil
}

func (s *sCommentUser) DeleteComment(
	ctx context.Context,
	command *commentCommand.DeleteCommentCommand,
) error {
	// 1. Find comment
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if commentFound == nil {
		return response.NewDataNotFoundError("comment not found")
	}

	// 2. Find post
	postFound, err := s.postRepo.GetById(ctx, commentFound.PostId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if postFound == nil {
		return response.NewDataNotFoundError("post not found")
	}

	// 3. Delete all child comment
	deletedCommentCount, err := s.commentRepo.DeleteCommentAndChildComment(ctx, commentFound.ID)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	updatePost := &postEntity.PostUpdate{
		CommentCount: pointer.Ptr(postFound.CommentCount - int(deletedCommentCount)),
	}

	err = updatePost.ValidatePostUpdate()
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePost)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if commentFound.ParentId == nil {
		return nil
	}

	// 5. Update rep_comment_count of parent comment -1
	parentCommentFound, err := s.commentRepo.GetById(ctx, *commentFound.ParentId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if parentCommentFound == nil {
		return response.NewDataNotFoundError("parent comment not found")
	}

	updateParentCommentData := commentEntity.CommentUpdate{
		RepCommentCount: pointer.Ptr(parentCommentFound.RepCommentCount - 1),
	}

	err = updateParentCommentData.ValidateCommentUpdate()
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.commentRepo.UpdateOne(ctx, parentCommentFound.ID, &updateParentCommentData)

	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 6. Delete comment report
	if err = s.commentReportRepor.DeleteByCommentId(ctx, command.CommentId); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sCommentUser) GetManyComments(
	ctx context.Context,
	query *commentQuery.GetManyCommentQuery,
) (result *commentQuery.GetManyCommentsResult, err error) {
	result = &commentQuery.GetManyCommentsResult{
		Comments:       nil,
		PagingResponse: nil,
	}

	// Get next layer of comment by root comment
	if query.ParentId != uuid.Nil {
		parentCommentFound, err := s.commentRepo.GetById(ctx, query.ParentId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		if parentCommentFound == nil {
			return nil, response.NewDataNotFoundError("parent comment not found")
		}

		queryResult, paging, err := s.commentRepo.GetMany(ctx, query)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		var commentResults []*common.CommentResultWithLiked
		for _, comment := range queryResult {
			commentResults = append(commentResults, mapper.NewCommentWithLikedResultFromEntity(comment))
		}

		result.Comments = commentResults
		result.PagingResponse = paging
		return result, nil
	} else {
		// Get first layer if it don't have parent id
		queryResult, paging, err := s.commentRepo.GetMany(ctx, query)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		var commentResults []*common.CommentResultWithLiked
		for _, comment := range queryResult {
			commentResults = append(commentResults, mapper.NewCommentWithLikedResultFromEntity(comment))
		}

		result.Comments = commentResults
		result.PagingResponse = paging
		return result, nil
	}
}
