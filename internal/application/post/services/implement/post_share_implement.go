package implement

import (
	"context"
	"errors"
	"github.com/google/uuid"
	post_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
)

type sPostShare struct {
	userRepo  post_repo.IUserRepository
	postRepo  post_repo.IPostRepository
	mediaRepo post_repo.IMediaRepository
}

func NewPostShareImplement(
	userRepo post_repo.IUserRepository,
	postRepo post_repo.IPostRepository,
	mediaRepo post_repo.IMediaRepository,
) *sPostShare {
	return &sPostShare{
		userRepo:  userRepo,
		postRepo:  postRepo,
		mediaRepo: mediaRepo,
	}
}

func (s *sPostShare) SharePost(
	ctx context.Context,
	postId uuid.UUID,
	userId uuid.UUID,
	shareInput *request.SharePostInput,
) (post *models.Post, resultCode int, httpStatusCode int, err error) {
	// 1. Find post by post_id
	postModel, err := s.postRepo.GetPost(ctx, "id = ?", postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, err
	}

	// 2. Create new post (parent_id = post_id, user_id = userId)
	if postModel.ParentId == nil {
		// 2.1. Copy post info from root post
		newPost := &models.Post{
			UserId:   userId,
			ParentId: &postModel.ID,
			Content:  shareInput.Content,
			Location: shareInput.Location,
			Privacy:  shareInput.Privacy,
		}

		// 2.2. Create new post
		newSharePost, err := s.postRepo.CreatePost(ctx, newPost)
		if err != nil {
			return nil, response.ErrServerFailed, http.StatusInternalServerError, err
		}

		return newSharePost, response.ErrCodeSuccess, http.StatusOK, nil
	} else {
		// 3. Find actually root post
		rootPost, err := s.postRepo.GetPost(ctx, "id=?", postModel.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, response.ErrDataNotFound, http.StatusBadRequest, err
			}
			return nil, response.ErrServerFailed, http.StatusInternalServerError, err
		}

		// 3.1. Copy post info from root post
		newPost := &models.Post{
			UserId:   userId,
			ParentId: &rootPost.ID,
			Content:  shareInput.Content,
			Location: shareInput.Location,
			Privacy:  shareInput.Privacy,
		}

		// 3.2. Create new post
		newSharePost, err := s.postRepo.CreatePost(ctx, newPost)
		if err != nil {
			return nil, response.ErrServerFailed, http.StatusInternalServerError, err
		}

		return newSharePost, response.ErrCodeSuccess, http.StatusOK, nil
	}
}
