package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	comment_command "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	comment_query "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	comment_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	comment_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
	"net/http"
)

type sCommentUser struct {
	commentRepo         comment_repo.ICommentRepository
	userRepo            comment_repo.IUserRepository
	postRepo            comment_repo.IPostRepository
	likeUserCommentRepo comment_repo.ILikeUserCommentRepository
}

func NewCommentUserImplement(
	commentRepo comment_repo.ICommentRepository,
	userRepo comment_repo.IUserRepository,
	postRepo comment_repo.IPostRepository,
	likeUserCommentRepo comment_repo.ILikeUserCommentRepository,
) *sCommentUser {
	return &sCommentUser{
		commentRepo:         commentRepo,
		userRepo:            userRepo,
		postRepo:            postRepo,
		likeUserCommentRepo: likeUserCommentRepo,
	}
}

func (s *sCommentUser) CreateComment(
	ctx context.Context,
	command *comment_command.CreateCommentCommand,
) (result *comment_command.CreateCommentResult, err error) {
	result = &comment_command.CreateCommentResult{}
	// 1. Find post
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
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
		return result, fmt.Errorf("Error when find post %w", err.Error())
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
			result.Comment = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("Error when find parent comment %w", err.Error())
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
			result.Comment = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("Error when update comment %w", err.Error())
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
			result.Comment = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("Error when update comment %w", err.Error())
		}

		// 2.4. Update rep count +1
		updateComment := &comment_entity.CommentUpdate{
			RepCommentCount: pointer.Ptr(parentComment.RepCommentCount + 1),
		}

		err = updateComment.ValidateCommentUpdate()
		if err != nil {
			result.Comment = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
		}

		_, err = s.commentRepo.UpdateOne(ctx, parentComment.ID, updateComment)

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
			return result, fmt.Errorf("Error when update comment %w", err.Error())
		}

		leftValue = rightValue
		rightValue++
	} else {
		// 3. Create comment if it don't have a parent id (root comment)
		maxRightValue, err := s.commentRepo.GetMaxCommentRightByPostId(ctx, command.PostId)
		if err != nil {
			result.Comment = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("Error when find max comment right: %w", err.Error())
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
	commentEntity, err := comment_entity.NewComment(
		command.PostId,
		command.UserId,
		command.ParentId,
		command.Content,
		leftValue,
		rightValue,
	)

	newComment, err := s.commentRepo.CreateOne(ctx, commentEntity)
	if err != nil {
		result.Comment = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when create comment %w", err.Error())
	}

	// 5. Update comment count for post
	updatePost := &post_entity.PostUpdate{
		CommentCount: pointer.Ptr(postFound.CommentCount + 1),
	}

	err = updatePost.ValidatePostUpdate()
	if err != nil {
		result.Comment = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update post %w", err.Error())
	}

	_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePost)
	if err != nil {
		result.Comment = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update post %w", err.Error())
	}

	result.Comment = mapper.NewCommentResultFromEntity(newComment)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentUser) UpdateComment(
	ctx context.Context,
	command *comment_command.UpdateCommentCommand,
) (result *comment_command.UpdateCommentResult, err error) {
	result = &comment_command.UpdateCommentResult{}
	updateData := &comment_entity.CommentUpdate{
		Content: command.Content,
	}

	err = updateData.ValidateCommentUpdate()
	if err != nil {
		result.Comment = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update comment %w", err.Error())
	}

	commentUpdate, err := s.commentRepo.UpdateOne(ctx, command.CommentId, updateData)
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
		return result, fmt.Errorf("Error when update comment %w", err.Error())
	}

	result.Comment = mapper.NewCommentResultFromEntity(commentUpdate)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentUser) DeleteComment(
	ctx context.Context,
	command *comment_command.DeleteCommentCommand,
) (result *comment_command.DeleteCommentResult, err error) {
	result = &comment_command.DeleteCommentResult{}
	// 1. Find comment
	fmt.Println(command.CommentId)
	comment, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when find comment %w", err.Error())
	}

	// 2. Find post
	postFound, err := s.postRepo.GetById(ctx, comment.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when find post %w", err.Error())
	}

	// 2. Define width to delete
	rightValue := comment.CommentRight
	leftValue := comment.CommentLeft
	width := rightValue - leftValue + 1

	// 3. Delete all child comment
	delete_conditions := map[string]interface{}{
		"post_id":         comment.PostId,
		"comment_left >=": leftValue,
		"comment_left <=": rightValue,
	}

	err = s.commentRepo.DeleteMany(ctx, delete_conditions)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update comment %w", err.Error())
	}

	// 4. Update rest of comment_right and comment_left
	update_conditions := map[string]interface{}{
		"post_id":        comment.PostId,
		"comment_left >": rightValue,
	}
	updateLeft := map[string]interface{}{
		"comment_left": gorm.Expr("comment_left - ?", width),
	}
	err = s.commentRepo.UpdateMany(ctx, update_conditions, updateLeft)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update comment %w", err.Error())
	}

	update_conditions = map[string]interface{}{
		"post_id":         comment.PostId,
		"comment_right >": rightValue,
	}
	update_right := map[string]interface{}{
		"comment_right": gorm.Expr("comment_right - ?", width),
	}
	err = s.commentRepo.UpdateMany(ctx, update_conditions, update_right)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update comment %w", err.Error())
	}

	updatePost := &post_entity.PostUpdate{
		CommentCount: pointer.Ptr(postFound.CommentCount - 1),
	}

	err = updatePost.ValidatePostUpdate()
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update post %w", err.Error())
	}

	_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePost)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update post %w", err.Error())
	}

	if comment.ParentId == nil {
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}

	// 5. Update rep_comment_count of parent comment -1
	parentComment, err := s.commentRepo.GetById(ctx, *comment.ParentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when find parent comment %w", err.Error())
	}

	updateParentCommentData := comment_entity.CommentUpdate{
		RepCommentCount: pointer.Ptr(parentComment.RepCommentCount - 1),
	}

	err = updateParentCommentData.ValidateCommentUpdate()
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update parent comment %w", err.Error())
	}

	_, err = s.commentRepo.UpdateOne(ctx, parentComment.ID, &updateParentCommentData)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("Error when update parent comment %w", err.Error())
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentUser) GetManyComments(
	ctx context.Context,
	query *comment_query.GetManyCommentQuery,
) (result *comment_query.GetManyCommentsResult, err error) {
	result = &comment_query.GetManyCommentsResult{}

	// Get next layer of comment by root comment
	if query.ParentId != uuid.Nil {
		_, err = s.commentRepo.GetById(ctx, query.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Comments = nil
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				result.PagingResponse = nil
				return result, err
			}
			result.Comments = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			result.PagingResponse = nil
			return result, fmt.Errorf("Error when find parent comment %w", err.Error())
		}

		queryResult, paging, err := s.commentRepo.GetMany(ctx, query)
		if err != nil {
			result.Comments = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			result.PagingResponse = nil
			return result, fmt.Errorf("Error when find parent comment %w", err.Error())
		}

		var commentResults []*common.CommentResultWithLiked
		for _, commentEntity := range queryResult {
			likeUserCommentEntity, err := comment_entity.NewLikeUserCommentEntity(query.AuthenticatedUserId, commentEntity.ID)
			if err != nil {
				result.Comments = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				result.PagingResponse = nil
				return result, fmt.Errorf("Error when find like user comment %w", err.Error())
			}

			isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)

			commentResults = append(commentResults, mapper.NewCommentWithLikedResultFromEntity(commentEntity, isLiked))
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
			result.Comments = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			result.PagingResponse = nil
			return result, fmt.Errorf("Error when find parent comment %w", err.Error())
		}

		var commentResults []*common.CommentResultWithLiked
		for _, commentEntity := range queryResult {
			likeUserCommentEntity, err := comment_entity.NewLikeUserCommentEntity(query.AuthenticatedUserId, commentEntity.ID)
			if err != nil {
				result.Comments = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				result.PagingResponse = nil
				return result, fmt.Errorf("Error when find like user comment %w", err.Error())
			}

			isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, likeUserCommentEntity)
			commentResults = append(commentResults, mapper.NewCommentWithLikedResultFromEntity(commentEntity, isLiked))
		}

		result.Comments = commentResults
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		result.PagingResponse = paging
		return result, nil
	}
}
