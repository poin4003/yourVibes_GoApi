package implement

import (
	"context"
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

func Test_GetPost_PostFoundInCache_PublicPrivacy(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()

	mockPostCache := new(mocks.MockPostCache)
	mockLikeUserPostRepo := new(mocks.MockLikeUserPostRepository)
	mockEventPublisher := newMockPostEventPublisher()

	service := &sPostUser{
		postCache:          mockPostCache,
		likeUserPostRepo:   mockLikeUserPostRepo,
		postEventPublisher: mockEventPublisher.PostEventPublisher,
	}

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
	mockEventPublisher.On("PublishStatistic", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	query := &postQuery.GetOnePostQuery{
		PostId:              postID,
		AuthenticatedUserId: userID,
	}

	result, err := service.GetPost(ctx, query)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockPost.ID, result.Post.ID)
}

func Test_GetPost_FriendOnlyAndIsFriend(t *testing.T) {
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
		postEventPublisher: mockEventPublisher.PostEventPublisher,
	}

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

	mockPostCache.On("GetPost", ctx, postID).Return(nil)
	mockPostRepo.On("GetById", ctx, postID).Return(mockPost, nil)
	mockFriendRepo.On("CheckFriendExist", ctx, &userEntity.Friend{
		UserId:   mockPost.UserId,
		FriendId: userID,
	}).Return(true, nil)
	mockPostCache.On("SetPost", ctx, mockPost).Return()
	mockLikeUserPostRepo.On("CheckUserLikePost", ctx, mock.Anything).Return(false, nil)
	mockEventPublisher.On("PublishStatistic", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	query := &postQuery.GetOnePostQuery{
		PostId:              postID,
		AuthenticatedUserId: userID,
	}

	result, err := service.GetPost(ctx, query)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockPost.ID, result.Post.ID)
}

func Test_GetPost_PostNotFoundInRepo(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()

	mockPostCache := new(mocks.MockPostCache)
	mockPostRepo := new(mocks.MockPostRepository)

	service := &sPostUser{
		postCache: mockPostCache,
		postRepo:  mockPostRepo,
	}

	mockPostCache.On("GetPost", ctx, postID).Return(nil)
	mockPostRepo.On("GetById", ctx, postID).Return(nil, nil)

	query := &postQuery.GetOnePostQuery{
		PostId:              postID,
		AuthenticatedUserId: userID,
	}

	result, err := service.GetPost(ctx, query)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "post not found")
}

func Test_GetPost_FriendOnly_NotFriend(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()

	response.InitCustomCode()

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
		postEventPublisher: mockEventPublisher.PostEventPublisher,
	}

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

	mockPostCache.On("GetPost", ctx, postID).Return(nil)
	mockPostRepo.On("GetById", ctx, postID).Return(mockPost, nil)
	mockFriendRepo.On("CheckFriendExist", ctx, &userEntity.Friend{
		UserId:   mockPost.UserId,
		FriendId: userID,
	}).Return(false, nil)
	mockPostCache.On("SetPost", ctx, mockPost).Return()
	mockLikeUserPostRepo.On("CheckUserLikePost", ctx, mock.Anything).Return(false, nil)
	mockEventPublisher.On("PublishStatistic", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	query := &postQuery.GetOnePostQuery{
		PostId:              postID,
		AuthenticatedUserId: userID,
	}

	result, err := service.GetPost(ctx, query)

	assert.Nil(t, result)
	assert.Error(t, err)

	expectedErr := response.NewCustomError(response.ErrPostFriendAccess)

	assert.Equal(t, expectedErr.Error(), err.Error())
}

func Test_GetPost_PrivatePrivacy_NotOwner(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()

	response.InitCustomCode()

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
		postEventPublisher: mockEventPublisher.PostEventPublisher,
	}

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

	mockPostCache.On("GetPost", ctx, postID).Return(nil)
	mockPostRepo.On("GetById", ctx, postID).Return(mockPost, nil)
	mockFriendRepo.On("CheckFriendExist", ctx, &userEntity.Friend{
		UserId:   mockPost.UserId,
		FriendId: userID,
	}).Return(false, nil)
	mockPostCache.On("SetPost", ctx, mockPost).Return()
	mockLikeUserPostRepo.On("CheckUserLikePost", ctx, mock.Anything).Return(false, nil)
	mockEventPublisher.On("PublishStatistic", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	query := &postQuery.GetOnePostQuery{
		PostId:              postID,
		AuthenticatedUserId: userID,
	}

	result, err := service.GetPost(ctx, query)
	assert.Nil(t, result)
	assert.Error(t, err)
	expectedErr := response.NewCustomError(response.ErrPostPrivateAccess)
	assert.Equal(t, expectedErr.Error(), err.Error())
}
