package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/stretchr/testify/mock"
)

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) GetById(ctx context.Context, id uuid.UUID) (*entities.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Post), args.Error(1)
}

func (m *MockPostRepository) CreateOne(ctx context.Context, entity *entities.Post) (*entities.Post, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Post), args.Error(1)
}

func (m *MockPostRepository) UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.PostUpdate) (*entities.Post, error) {
	args := m.Called(ctx, id, updateData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Post), args.Error(1)
}

func (m *MockPostRepository) UpdateMany(ctx context.Context, condition map[string]interface{}, updateData *entities.PostUpdate) error {
	args := m.Called(ctx, condition, updateData)
	return args.Error(0)
}

func (m *MockPostRepository) DeleteOne(ctx context.Context, id uuid.UUID) (*entities.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Post), args.Error(1)
}

func (m *MockPostRepository) GetMany(ctx context.Context, query *query.GetManyPostQuery) ([]*entities.Post, *response.PagingResponse, error) {
	args := m.Called(ctx, query)
	posts, _ := args.Get(0).([]*entities.Post)
	paging, _ := args.Get(1).(*response.PagingResponse)
	return posts, paging, args.Error(2)
}

func (m *MockPostRepository) GetTrendingPost(ctx context.Context, query *query.GetTrendingPostQuery) ([]*entities.Post, *response.PagingResponse, error) {
	args := m.Called(ctx, query)
	posts, _ := args.Get(0).([]*entities.Post)
	paging, _ := args.Get(1).(*response.PagingResponse)
	return posts, paging, args.Error(2)
}

func (m *MockPostRepository) UpdateExpiredAdvertisements(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPostRepository) CheckPostOwner(ctx context.Context, postId uuid.UUID, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, postId, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockPostRepository) GetTotalPostCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockPostRepository) GetTotalPostCountByUserId(ctx context.Context, userId uuid.UUID) (int64, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostRepository) UpdatePostAndStatistics(ctx context.Context, postID uuid.UUID, likeDelta, commentDelta, clicksDelta, impressionDelta, reachDelta int) error {
	args := m.Called(ctx, postID, likeDelta, commentDelta, clicksDelta, impressionDelta, reachDelta)
	return args.Error(0)
}

func (m *MockPostRepository) DelayPostCreatedAt(ctx context.Context, postID uuid.UUID, delay time.Duration) error {
	args := m.Called(ctx, postID, delay)
	return args.Error(0)
}
