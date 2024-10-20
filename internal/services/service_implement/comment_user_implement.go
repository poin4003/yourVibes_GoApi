package service_implement

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type sCommentUser struct {
	commentRepo repository.ICommentRepository
	userRepo    repository.IUserRepository
	postRepo    repository.IPostRepository
}

func NewCommentUserImplement(
	commentRepo repository.ICommentRepository,
	userRepo repository.IUserRepository,
	postRepo repository.IPostRepository,
) *sCommentUser {
	return &sCommentUser{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		postRepo:    postRepo,
	}
}

func (s *sCommentUser) CreateComment(
	ctx context.Context,
	commentModel *model.Comment,
) (comment *model.Comment, resultCode int, err error) {
	post, err := s.postRepo.GetPost(ctx, "id=?", commentModel.PostId)
	if err != nil {
		return nil, response.ErrServerFailed, fmt.Errorf("Error when find post %w", err.Error())
	}

	if post == nil {
		return nil, response.ErrServerFailed, fmt.Errorf("Post not found")
	}

	var rightValue int

	if commentModel.ParentId != nil {
		parentComment, err := s.commentRepo.GetOneComment(ctx, "id=?", *commentModel.ParentId)
		if err != nil {
			return nil, response.ErrServerFailed, fmt.Errorf("Error when find parent comment %w", err.Error())
		}

		if parentComment == nil {
			return nil, response.ErrDataNotFound, fmt.Errorf("Parent comment not found")
		}

		rightValue = parentComment.CommentRight

		// Find comment by postId and update all comment.comment_right +2 if that comment.comment_right greater than or equal rightValue
		conditions := map[string]interface{}{
			"post_id":          commentModel.PostId,
			"comment_right >=": rightValue,
		}
		updateRight := map[string]interface{}{
			"comment_right": gorm.Expr("comment_right + ?", 2),
		}
		err = s.commentRepo.UpdateManyComment(ctx, conditions, updateRight)
		if err != nil {
			return nil, response.ErrServerFailed, fmt.Errorf("Error when update comment %w", err.Error())
		}

		// Find comment by postId and update all comment.comment_left +2 if that comment.comment_left greater than rightValue
		conditions = map[string]interface{}{
			"post_id":        commentModel.PostId,
			"comment_left >": rightValue,
		}
		updateLeft := map[string]interface{}{
			"comment_left": gorm.Expr("comment_left + ?", 2),
		}
		err = s.commentRepo.UpdateManyComment(ctx, conditions, updateLeft)
		if err != nil {
			return nil, response.ErrServerFailed, fmt.Errorf("Error when update comment %w", err.Error())
		}

		// Update rep count +1
		parentComment.RepCommentCount++
		_, err = s.commentRepo.UpdateOneComment(ctx, parentComment.ID, map[string]interface{}{
			"rep_comment_count": parentComment.RepCommentCount,
		})

		commentModel.CommentLeft = rightValue
		commentModel.CommentRight = rightValue + 1
	} else {
		maxRightValue, err := s.commentRepo.GetMaxCommentRightByPostId(ctx, commentModel.PostId)
		if err != nil {
			return nil, response.ErrServerFailed, fmt.Errorf("Error when find max comment right: %w", err.Error())
		}

		if maxRightValue != 0 {
			rightValue = maxRightValue + 1
		} else {
			rightValue = 1
		}

		commentModel.CommentLeft = rightValue
		commentModel.CommentRight = rightValue + 1
	}

	newComment, err := s.commentRepo.CreateComment(ctx, commentModel)
	if err != nil {
		return nil, response.ErrServerFailed, fmt.Errorf("Error when create comment %w", err.Error())
	}

	return newComment, response.ErrCodeSuccess, nil
}

func (s *sCommentUser) UpdateComment(
	ctx context.Context,
	commentId uuid.UUID,
	updateData map[string]interface{},
) (comment *model.Comment, resultCode int, err error) {
	commentUpdate, err := s.commentRepo.UpdateOneComment(ctx, commentId, updateData)
	if err != nil {
		return nil, response.ErrServerFailed, fmt.Errorf("Error when update comment %w", err.Error())
	}

	return commentUpdate, response.ErrCodeSuccess, nil
}

func (s *sCommentUser) DeleteComment(
	ctx context.Context,
	commentId uuid.UUID,
) (resultCode int, err error) {
	// 1. Find comment
	comment, err := s.commentRepo.GetOneComment(ctx, "id=?", commentId)
	if err != nil {
		return response.ErrServerFailed, fmt.Errorf("Error when find comment %w", err.Error())
	}

	if comment == nil {
		return response.ErrDataNotFound, fmt.Errorf("Comment not found")
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

	err = s.commentRepo.DeleteManyComment(ctx, delete_conditions)
	if err != nil {
		return response.ErrServerFailed, fmt.Errorf("Error when update comment %w", err.Error())
	}

	// 4. Update rest of comment_right and comment_left
	update_conditions := map[string]interface{}{
		"post_id":        comment.PostId,
		"comment_left >": rightValue,
	}
	updateLeft := map[string]interface{}{
		"comment_left": gorm.Expr("comment_left - ?", width),
	}
	err = s.commentRepo.UpdateManyComment(ctx, update_conditions, updateLeft)
	if err != nil {
		return response.ErrServerFailed, fmt.Errorf("Error when update comment %w", err.Error())
	}

	update_conditions = map[string]interface{}{
		"post_id":         comment.PostId,
		"comment_right >": rightValue,
	}
	update_right := map[string]interface{}{
		"comment_right": gorm.Expr("comment_right - ?", width),
	}
	err = s.commentRepo.UpdateManyComment(ctx, update_conditions, update_right)
	if err != nil {
		return response.ErrServerFailed, fmt.Errorf("Error when update comment %w", err.Error())
	}

	// 5. Update rep_comment_count of parent comment -1
	parentComment, err := s.commentRepo.GetOneComment(ctx, "id=?", comment.ParentId)
	if err != nil {
		return response.ErrServerFailed, fmt.Errorf("Error when find parent comment %w", err.Error())
	}

	if parentComment == nil {
		return response.ErrDataNotFound, fmt.Errorf("Parent comment not found")
	}

	parentComment.RepCommentCount--
	_, err = s.commentRepo.UpdateOneComment(ctx, parentComment.ID, map[string]interface{}{
		"rep_comment_count": parentComment.RepCommentCount,
	})

	return response.ErrCodeSuccess, nil
}

func (s *sCommentUser) GetManyComments(
	ctx context.Context,
	query *query_object.CommentQueryObject,
) (comments []*model.Comment, resultCode int, paingResponse *response.PagingResponse, err error) {
	var queryResult []*model.Comment

	if query.ParentId != "" {
		parentComment, err := s.commentRepo.GetOneComment(ctx, "id=?", query.ParentId)
		if err != nil {
			return nil, response.ErrServerFailed, nil, fmt.Errorf("Error when find parent comment %w", err.Error())
		}
		if parentComment == nil {
			return nil, response.ErrDataNotFound, nil, fmt.Errorf("Parent comment not found")
		}

		queryResult, pagingResponse, err := s.commentRepo.GetManyComment(ctx, query)
		if err != nil {
			return nil, response.ErrServerFailed, nil, fmt.Errorf("Error when find parent comment %w", err.Error())
		}

		return queryResult, response.ErrCodeSuccess, pagingResponse, nil
	} else {
		queryResult, paingResponse, err = s.commentRepo.GetManyComment(ctx, query)
		if err != nil {
			return nil, response.ErrServerFailed, nil, fmt.Errorf("Error when find parent comment %w", err.Error())
		}
	}

	return queryResult, response.ErrCodeSuccess, paingResponse, nil
}
