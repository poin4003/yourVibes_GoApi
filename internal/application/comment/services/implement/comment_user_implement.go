package implement

import (
	"context"
	"errors"
	"fmt"
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
	"gorm.io/gorm"
	"net/http"
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
	result = &commentCommand.CreateCommentResult{
		Comment:        nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Find post
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Comment = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("error when find post %v", err.Error())
	}

	// 2. Defined value for nested set model
	var rightValue, leftValue int

	if command.ParentId != nil {
		// 2.1. Get root comment
		parentComment, err := s.commentRepo.GetById(ctx, *command.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Comment = nil
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			return result, fmt.Errorf("error when find parent comment %v", err.Error())
		}

		rightValue = parentComment.CommentRight

		// 2.2. Find comment by postId and update all comment.comment_right +2 if that comment.comment_right greater than or equal rightValue
		conditions := map[string]interface{}{
			"post_id":          command.PostId,
			"comment_right >=": rightValue,
		}
		updateRight := map[string]interface{}{
			"comment_right": gorm.Expr("comment_right + ?", 2),
		}
		err = s.commentRepo.UpdateMany(ctx, conditions, updateRight)
		if err != nil {
			return result, fmt.Errorf("error when update comment %v", err.Error())
		}

		// 2.3. Find comment by postId and update all comment.comment_left +2 if that comment.comment_left greater than rightValue
		conditions = map[string]interface{}{
			"post_id":        command.PostId,
			"comment_left >": rightValue,
		}
		updateLeft := map[string]interface{}{
			"comment_left": gorm.Expr("comment_left + ?", 2),
		}
		err = s.commentRepo.UpdateMany(ctx, conditions, updateLeft)
		if err != nil {
			return result, fmt.Errorf("error when update comment %v", err.Error())
		}

		// 2.4. Update rep count +1
		updateComment := &commentEntity.CommentUpdate{
			RepCommentCount: pointer.Ptr(parentComment.RepCommentCount + 1),
		}

		if err = updateComment.ValidateCommentUpdate(); err != nil {
			return result, err
		}

		_, err = s.commentRepo.UpdateOne(ctx, parentComment.ID, updateComment)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Comment = nil
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			return result, fmt.Errorf("error when update comment %v", err.Error())
		}

		leftValue = rightValue
		rightValue++
	} else {
		// 3. Create comment if it don't have a parent id (root comment)
		maxRightValue, err := s.commentRepo.GetMaxCommentRightByPostId(ctx, command.PostId)
		if err != nil {
			return result, fmt.Errorf("error when find max comment right: %v", err.Error())
		}

		// 3.1. Assign default value if it a root comment
		if maxRightValue != 0 {
			rightValue = maxRightValue + 1
		} else {
			rightValue = 1
		}

		leftValue = rightValue
		rightValue++
	}

	// 4. Create a comment
	newComment, err := commentEntity.NewComment(
		command.PostId,
		command.UserId,
		command.ParentId,
		command.Content,
		leftValue,
		rightValue,
	)

	commentCreated, err := s.commentRepo.CreateOne(ctx, newComment)
	if err != nil {
		return result, fmt.Errorf("error when create comment %v", err.Error())
	}

	// 5. Update comment count for post
	updatePost := &postEntity.PostUpdate{
		CommentCount: pointer.Ptr(postFound.CommentCount + 1),
	}

	err = updatePost.ValidatePostUpdate()
	if err != nil {
		return result, fmt.Errorf("error when update post %v", err.Error())
	}

	_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePost)
	if err != nil {
		return result, fmt.Errorf("error when update post %v", err.Error())
	}

	// 6. Validate comment after create
	validateComment, err := commentValidator.NewValidatedComment(commentCreated)
	if err != nil {
		return result, fmt.Errorf("failed to validate comment: %w", err)
	}

	result.Comment = mapper.NewCommentResultFromValidateEntity(validateComment)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentUser) UpdateComment(
	ctx context.Context,
	command *commentCommand.UpdateCommentCommand,
) (result *commentCommand.UpdateCommentResult, err error) {
	result = &commentCommand.UpdateCommentResult{
		Comment:        nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	updateData := &commentEntity.CommentUpdate{
		Content: command.Content,
	}

	err = updateData.ValidateCommentUpdate()
	if err != nil {
		return result, fmt.Errorf("error when update comment %v", err.Error())
	}

	commentUpdate, err := s.commentRepo.UpdateOne(ctx, command.CommentId, updateData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Comment = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("error when update comment %v", err.Error())
	}

	result.Comment = mapper.NewCommentResultFromEntity(commentUpdate)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentUser) DeleteComment(
	ctx context.Context,
	command *commentCommand.DeleteCommentCommand,
) (result *commentCommand.DeleteCommentResult, err error) {
	result = &commentCommand.DeleteCommentResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Find comment
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("error when find comment %v", err.Error())
	}

	// 2. Find post
	postFound, err := s.postRepo.GetById(ctx, commentFound.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("error when find post %v", err.Error())
	}

	// 2. Define width to delete
	rightValue := commentFound.CommentRight
	leftValue := commentFound.CommentLeft
	width := rightValue - leftValue + 1

	// 3. Delete all child comment
	deleteConditions := map[string]interface{}{
		"post_id":         commentFound.PostId,
		"comment_left >=": leftValue,
		"comment_left <=": rightValue,
	}

	deletedCommentCount, err := s.commentRepo.DeleteMany(ctx, deleteConditions)
	if err != nil {
		return result, fmt.Errorf("error when update comment %v", err.Error())
	}

	// 4. Update rest of comment_right and comment_left
	updateConditions := map[string]interface{}{
		"post_id":        commentFound.PostId,
		"comment_left >": rightValue,
	}
	updateLeft := map[string]interface{}{
		"comment_left": gorm.Expr("comment_left - ?", width),
	}
	err = s.commentRepo.UpdateMany(ctx, updateConditions, updateLeft)
	if err != nil {
		return result, fmt.Errorf("error when update comment %v", err.Error())
	}

	updateConditions = map[string]interface{}{
		"post_id":         commentFound.PostId,
		"comment_right >": rightValue,
	}
	updateRight := map[string]interface{}{
		"comment_right": gorm.Expr("comment_right - ?", width),
	}
	err = s.commentRepo.UpdateMany(ctx, updateConditions, updateRight)
	if err != nil {
		return result, fmt.Errorf("error when update comment %v", err.Error())
	}

	updatePost := &postEntity.PostUpdate{
		CommentCount: pointer.Ptr(postFound.CommentCount - int(deletedCommentCount)),
	}

	err = updatePost.ValidatePostUpdate()
	if err != nil {
		return result, fmt.Errorf("error when update post %v", err.Error())
	}

	_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePost)
	if err != nil {
		return result, fmt.Errorf("error when update post %v", err.Error())
	}

	if commentFound.ParentId == nil {
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}

	// 5. Update rep_comment_count of parent comment -1
	parentCommentFound, err := s.commentRepo.GetById(ctx, *commentFound.ParentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("error when find parent comment %v", err.Error())
	}

	updateParentCommentData := commentEntity.CommentUpdate{
		RepCommentCount: pointer.Ptr(parentCommentFound.RepCommentCount - 1),
	}

	err = updateParentCommentData.ValidateCommentUpdate()
	if err != nil {
		return result, fmt.Errorf("error when update parent comment %v", err.Error())
	}

	_, err = s.commentRepo.UpdateOne(ctx, parentCommentFound.ID, &updateParentCommentData)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("error when update parent comment %v", err.Error())
	}

	// 6. Delete comment report
	if err = s.commentReportRepor.DeleteByCommentId(ctx, command.CommentId); err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentUser) GetManyComments(
	ctx context.Context,
	query *commentQuery.GetManyCommentQuery,
) (result *commentQuery.GetManyCommentsResult, err error) {
	result = &commentQuery.GetManyCommentsResult{
		Comments:       nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
		PagingResponse: nil,
	}

	// Get next layer of comment by root comment
	if query.ParentId != uuid.Nil {
		_, err = s.commentRepo.GetById(ctx, query.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			return result, fmt.Errorf("error when find parent comment %v", err.Error())
		}

		queryResult, paging, err := s.commentRepo.GetMany(ctx, query)
		if err != nil {
			return result, fmt.Errorf("error when find parent comment %v", err.Error())
		}

		var commentResults []*common.CommentResultWithLiked
		for _, comment := range queryResult {
			commentResults = append(commentResults, mapper.NewCommentWithLikedResultFromEntity(comment))
		}

		result.Comments = commentResults
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		result.PagingResponse = paging
		return result, nil
	} else {
		// Get first layer if it don't have parent id
		queryResult, paging, err := s.commentRepo.GetMany(ctx, query)
		if err != nil {
			return result, fmt.Errorf("error when find parent comment %v", err.Error())
		}

		var commentResults []*common.CommentResultWithLiked
		for _, comment := range queryResult {
			commentResults = append(commentResults, mapper.NewCommentWithLikedResultFromEntity(comment))
		}

		result.Comments = commentResults
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		result.PagingResponse = paging
		return result, nil
	}
}
