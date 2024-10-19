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
		parentComment, err := s.commentRepo.GetComment(ctx, "id=?", *commentModel.ParentId)
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
	return &model.Comment{}, 0, nil
}

func (s *sCommentUser) DeleteComment(
	ctx context.Context,
	commentId uuid.UUID,
) (resultCode int, err error) {
	return 0, nil
}

func (s *sCommentUser) GetManyComments(
	ctx context.Context,
	query *query_object.CommentQueryObject,
) (comments []*model.Comment, resultCode int, err error) {
	var queryResult []*model.Comment

	if query.ParentId != "" {
		parentComment, err := s.commentRepo.GetComment(ctx, "id=?", query.ParentId)
		if err != nil {
			return nil, response.ErrServerFailed, fmt.Errorf("Error when find parent comment %w", err.Error())
		}
		if parentComment == nil {
			return nil, response.ErrDataNotFound, fmt.Errorf("Parent comment not found")
		}

		queryResult, err := s.commentRepo.GetManyComment(ctx, query)
		if err != nil {
			return nil, response.ErrServerFailed, fmt.Errorf("Error when find parent comment %w", err.Error())
		}

		return queryResult, response.ErrCodeSuccess, nil
	} else {
		queryResult, err = s.commentRepo.GetManyComment(ctx, query)
		if err != nil {
			return nil, response.ErrServerFailed, fmt.Errorf("Error when find parent comment %w", err.Error())
		}
	}

	return queryResult, response.ErrCodeSuccess, nil
}
