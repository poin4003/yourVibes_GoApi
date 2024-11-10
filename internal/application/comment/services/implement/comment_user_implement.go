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
	commentModel *models.Comment,
) (comment *models.Comment, resultCode int, httpStatusCode int, err error) {
	// 1. Find post
	postFound, err := s.postRepo.GetPost(ctx, "id=?", commentModel.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg_response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find post %w", err.Error())
	}

	var rightValue int

	if commentModel.ParentId != nil {
		parentComment, err := s.commentRepo.GetOneComment(ctx, "id=?", *commentModel.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg_response.ErrDataNotFound, http.StatusBadRequest, err
			}
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find parent comment %w", err.Error())
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
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when update comment %w", err.Error())
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
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when update comment %w", err.Error())
		}

		// Update rep count +1
		parentComment.RepCommentCount++
		_, err = s.commentRepo.UpdateOneComment(ctx, parentComment.ID, map[string]interface{}{
			"rep_comment_count": parentComment.RepCommentCount,
		})

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg_response.ErrDataNotFound, http.StatusInternalServerError, err
			}
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when update comment %w", err.Error())
		}

		commentModel.CommentLeft = rightValue
		commentModel.CommentRight = rightValue + 1
	} else {
		maxRightValue, err := s.commentRepo.GetMaxCommentRightByPostId(ctx, commentModel.PostId)
		if err != nil {
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find max comment right: %w", err.Error())
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
		return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when create comment %w", err.Error())
	}

	postFound.CommentCount++
	_, err = s.postRepo.UpdatePost(ctx, postFound.ID, map[string]interface{}{
		"comment_count": postFound.CommentCount,
	})

	return newComment, pkg_response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sCommentUser) UpdateComment(
	ctx context.Context,
	commentId uuid.UUID,
	updateData map[string]interface{},
) (comment *models.Comment, resultCode int, httpStatusCode int, err error) {
	commentUpdate, err := s.commentRepo.UpdateOneComment(ctx, commentId, updateData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg_response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when update comment %w", err.Error())
	}

	return commentUpdate, pkg_response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sCommentUser) DeleteComment(
	ctx context.Context,
	commentId uuid.UUID,
) (resultCode int, httpStatusCode int, err error) {
	// 1. Find comment
	comment, err := s.commentRepo.GetOneComment(ctx, "id=?", commentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg_response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find comment %w", err.Error())
	}

	// 2. Find post
	postFound, err := s.postRepo.GetPost(ctx, "id=?", comment.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg_response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find post %w", err.Error())
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
		return pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when update comment %w", err.Error())
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
		return pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when update comment %w", err.Error())
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
		return pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when update comment %w", err.Error())
	}

	postFound.CommentCount--
	_, err = s.postRepo.UpdatePost(ctx, postFound.ID, map[string]interface{}{
		"comment_count": postFound.CommentCount,
	})

	if err != nil {
		return pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when update comment count %w", err.Error())
	}

	if comment.ParentComment == nil {
		return pkg_response.ErrCodeSuccess, http.StatusOK, nil
	}

	// 5. Update rep_comment_count of parent comment -1
	parentComment, err := s.commentRepo.GetOneComment(ctx, "id=?", comment.ParentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg_response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find parent comment %w", err.Error())
	}

	parentComment.RepCommentCount--
	_, err = s.commentRepo.UpdateOneComment(ctx, parentComment.ID, map[string]interface{}{
		"rep_comment_count": parentComment.RepCommentCount,
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg_response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return pkg_response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when update parent comment %w", err.Error())
	}

	return pkg_response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sCommentUser) GetManyComments(
	ctx context.Context,
	query *query.CommentQueryObject,
	userId uuid.UUID,
) (commentDtos []*response.CommentDto, resultCode int, httpStatusCode int, paingpkg_response *pkg_response.PagingResponse, err error) {
	var queryResult []*models.Comment

	if query.ParentId != "" {
		_, err = s.commentRepo.GetOneComment(ctx, "id=?", query.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg_response.ErrDataNotFound, http.StatusNotFound, nil, err
			}
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, nil, fmt.Errorf("Error when find parent comment %w", err.Error())
		}

		queryResult, pagingpkg_response, err := s.commentRepo.GetManyComment(ctx, query)
		if err != nil {
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, nil, fmt.Errorf("Error when find parent comment %w", err.Error())
		}

		for _, comment := range queryResult {
			isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, &models.LikeUserComment{
				CommentId: comment.ID,
				UserId:    userId,
			})

			commentDto := mapper.MapCommentToCommentDto(comment, isLiked)
			commentDtos = append(commentDtos, commentDto)
		}

		return commentDtos, pkg_response.ErrCodeSuccess, http.StatusOK, pagingpkg_response, nil
	} else {
		queryResult, paingpkg_response, err = s.commentRepo.GetManyComment(ctx, query)
		if err != nil {
			return nil, pkg_response.ErrServerFailed, http.StatusInternalServerError, nil, fmt.Errorf("Error when find parent comment %w", err.Error())
		}

		for _, comment := range queryResult {
			isLiked, _ := s.likeUserCommentRepo.CheckUserLikeComment(ctx, &models.LikeUserComment{
				CommentId: comment.ID,
				UserId:    userId,
			})

			commentDto := mapper.MapCommentToCommentDto(comment, isLiked)
			commentDtos = append(commentDtos, commentDto)
		}

		return commentDtos, pkg_response.ErrCodeSuccess, http.StatusOK, paingpkg_response, nil
	}
}
