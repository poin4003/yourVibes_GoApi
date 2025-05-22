package implement

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/producer"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services/mocks"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPostEventPublisher struct {
	*producer.PostEventPublisher
	mock.Mock
}

func newMockPostEventPublisher() *mockPostEventPublisher {
	pub := new(producer.PostEventPublisher) // khởi tạo struct rỗng, không set trường nào
	return &mockPostEventPublisher{PostEventPublisher: pub}
}

func (m *mockPostEventPublisher) PublishStatistic(ctx context.Context, msg interface{}, routingKey string) error {
	args := m.Called(ctx, msg, routingKey)
	return args.Error(0)
}

func Test_sPostUser_GetPost(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()

	mockPostCache := new(mocks.MockPostCache)
	mockPostRepo := new(mocks.MockPostRepository)
	mockFriendRepo := new(mocks.MockFriendRepository)
	mockLikeUserPostRepo := new(mocks.MockLikeUserPostRepository)
	mockEventPublisher := newMockPostEventPublisher()

	service := &sPostUser{
		postCache:          mockPostCache,
		postRepo:           mockPostRepo,
		friendRepo:         mockFriendRepo,
		likeUserPostRepo:   mockLikeUserPostRepo,
		postEventPublisher: mockEventPublisher.PostEventPublisher, // truyền đúng kiểu
	}

	query := &postQuery.GetOnePostQuery{
		PostId:              postID,
		AuthenticatedUserId: userID,
	}

	t.Run("post found in cache, public privacy", func(t *testing.T) {
		defer mockPostCache.AssertExpectations(t)
		defer mockLikeUserPostRepo.AssertExpectations(t)
		defer mockEventPublisher.AssertExpectations(t)

		mockPost := &postEntity.Post{
			ID:      postID,
			UserId:  uuid.New(),
			Privacy: consts.PUBLIC,
			User: &postEntity.User{
				ID:         uuid.New(),
				Name:       "test",
				FamilyName: "test",
				AvatarUrl:  "test",
			},
			Content: "test",
		}

		mockPostCache.On("GetPost", ctx, postID).Return(mockPost)
		mockLikeUserPostRepo.On("CheckUserLikePost", ctx, mock.Anything).Return(true, nil)
		mockEventPublisher.
			On("PublishStatistic", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Maybe()

		result, err := service.GetPost(ctx, query)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockPost.ID, result.Post.ID)
	})

	t.Run("post not found in cache, found in repo, friend only privacy and is friend", func(t *testing.T) {
		defer mockPostCache.AssertExpectations(t)
		defer mockPostRepo.AssertExpectations(t)
		defer mockFriendRepo.AssertExpectations(t)
		defer mockLikeUserPostRepo.AssertExpectations(t)
		defer mockEventPublisher.AssertExpectations(t)

		// Reset expectations để tránh ảnh hưởng test khác
		mockPostCache.ExpectedCalls = nil
		mockPostRepo.ExpectedCalls = nil
		mockFriendRepo.ExpectedCalls = nil
		mockLikeUserPostRepo.ExpectedCalls = nil
		mockEventPublisher.ExpectedCalls = nil

		mockPostCache.On("GetPost", ctx, postID).Return(nil)
		mockPost := &postEntity.Post{
			ID:      postID,
			UserId:  uuid.New(),
			Privacy: consts.FRIEND_ONLY,
			User: &postEntity.User{
				ID:         uuid.New(),
				Name:       "test",
				FamilyName: "test",
				AvatarUrl:  "test",
			},
			Content: "test",
		}
		mockPostRepo.On("GetById", ctx, postID).Return(mockPost, nil)
		mockFriendRepo.On("CheckFriendExist", ctx, &userEntity.Friend{
			UserId:   mockPost.UserId,
			FriendId: userID,
		}).Return(true, nil)
		mockPostCache.On("SetPost", ctx, mockPost).Return()
		mockLikeUserPostRepo.On("CheckUserLikePost", ctx, mock.Anything).Return(false, nil)
		mockEventPublisher.
			On("PublishStatistic", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Maybe()

		result, err := service.GetPost(ctx, query)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockPost.ID, result.Post.ID)
	})

	t.Run("post not found in repo", func(t *testing.T) {
		defer mockPostCache.AssertExpectations(t)
		defer mockPostRepo.AssertExpectations(t)

		mockPostCache.ExpectedCalls = nil
		mockPostRepo.ExpectedCalls = nil

		mockPostCache.On("GetPost", ctx, postID).Return(nil)
		mockPostRepo.On("GetById", ctx, postID).Return(nil, nil)

		result, err := service.GetPost(ctx, query)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post not found")
	})

	t.Run("friend only privacy, not a friend", func(t *testing.T) {
		defer mockPostCache.AssertExpectations(t)
		defer mockPostRepo.AssertExpectations(t)
		defer mockFriendRepo.AssertExpectations(t)

		mockPostCache.ExpectedCalls = nil
		mockPostRepo.ExpectedCalls = nil
		mockFriendRepo.ExpectedCalls = nil

		mockPostCache.On("GetPost", ctx, postID).Return(nil)
		mockPostCache.On("SetPost", mock.Anything, mock.Anything).Return()

		mockPost := &postEntity.Post{
			ID:      postID,
			UserId:  uuid.New(),
			Privacy: consts.FRIEND_ONLY,
			User: &postEntity.User{
				ID:         uuid.New(),
				Name:       "test",
				FamilyName: "test",
				AvatarUrl:  "test",
			},
			Content: "test",
		}
		mockPostRepo.On("GetById", ctx, postID).Return(mockPost, nil)
		mockFriendRepo.On("CheckFriendExist", ctx, &userEntity.Friend{
			UserId:   mockPost.UserId,
			FriendId: userID,
		}).Return(false, nil)

		result, err := service.GetPost(ctx, query)

		if err != nil {
			fmt.Println("Error Friend", err.Error())
		}
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, response.ErrPostFriendAccess, err)

	})

	t.Run("private privacy, not owner", func(t *testing.T) {
		defer mockPostCache.AssertExpectations(t)
		defer mockPostRepo.AssertExpectations(t)

		mockPostCache.ExpectedCalls = nil
		mockPostRepo.ExpectedCalls = nil

		mockPostCache.On("GetPost", ctx, postID).Return(nil)
		mockPost := &postEntity.Post{
			ID:      postID,
			UserId:  uuid.New(),
			Privacy: consts.PRIVATE,
			User: &postEntity.User{
				ID:         uuid.New(),
				Name:       "test",
				FamilyName: "test",
				AvatarUrl:  "test",
			},
			Content: "test",
		}
		mockPostRepo.On("GetById", ctx, postID).Return(mockPost, nil)

		result, err := service.GetPost(ctx, query)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), response.ErrPostPrivateAccess)
	})
}
