package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/stretchr/testify/mock"
)

type MockPostCache struct {
	mock.Mock
}

func (m *MockPostCache) SetPost(ctx context.Context, post *entities.Post) {
	m.Called(ctx, post)
}

func (m *MockPostCache) GetPost(ctx context.Context, postID uuid.UUID) *entities.Post {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*entities.Post)
}

func (m *MockPostCache) DeletePost(ctx context.Context, postID uuid.UUID) {
	m.Called(ctx, postID)
}

func (m *MockPostCache) SetFeeds(ctx context.Context, inputKey consts.RedisKey, userID uuid.UUID, postsIds []uuid.UUID, paging *response.PagingResponse) {
	m.Called(ctx, inputKey, userID, postsIds, paging)
}

func (m *MockPostCache) GetFeeds(ctx context.Context, inputKey consts.RedisKey, userID uuid.UUID, limit, page int) ([]uuid.UUID, *response.PagingResponse) {
	args := m.Called(ctx, inputKey, userID, limit, page)
	ids, _ := args.Get(0).([]uuid.UUID)
	paging, _ := args.Get(1).(*response.PagingResponse)
	return ids, paging
}

func (m *MockPostCache) DeleteFeeds(ctx context.Context, inputKey consts.RedisKey, userID uuid.UUID) {
	m.Called(ctx, inputKey, userID)
}

func (m *MockPostCache) DeleteFriendFeeds(ctx context.Context, inputKey consts.RedisKey, friendIDs []uuid.UUID) {
	m.Called(ctx, inputKey, friendIDs)
}

func (m *MockPostCache) DeleteRelatedPost(ctx context.Context, inputKey consts.RedisKey, userID uuid.UUID) {
	m.Called(ctx, inputKey, userID)
}

func (m *MockPostCache) SetPostForCreate(ctx context.Context, postID uuid.UUID, post *entities.PostForCreate) error {
	args := m.Called(ctx, postID, post)
	return args.Error(0)
}

func (m *MockPostCache) GetPostForCreate(ctx context.Context, postID uuid.UUID) (*entities.PostForCreate, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PostForCreate), args.Error(1)
}

func (m *MockPostCache) DeletePostForCreate(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockPostCache) DeleteAllPostCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
