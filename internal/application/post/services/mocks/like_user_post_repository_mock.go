package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/stretchr/testify/mock"
)

type MockLikeUserPostRepository struct {
	mock.Mock
}

func (m *MockLikeUserPostRepository) CreateLikeUserPost(ctx context.Context, entity *entities.LikeUserPost) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockLikeUserPostRepository) DeleteLikeUserPost(ctx context.Context, entity *entities.LikeUserPost) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockLikeUserPostRepository) GetLikeUserPost(ctx context.Context, query *query.GetPostLikeQuery) ([]*entities.User, *response.PagingResponse, error) {
	args := m.Called(ctx, query)
	users, _ := args.Get(0).([]*entities.User)
	paging, _ := args.Get(1).(*response.PagingResponse)
	return users, paging, args.Error(2)
}

func (m *MockLikeUserPostRepository) CheckUserLikePost(ctx context.Context, entity *entities.LikeUserPost) (bool, error) {
	args := m.Called(ctx, entity)
	return args.Bool(0), args.Error(1)
}

func (m *MockLikeUserPostRepository) CheckUserLikeManyPost(ctx context.Context, query *query.CheckUserLikeManyPostQuery) (map[uuid.UUID]bool, error) {
	args := m.Called(ctx, query)
	res, _ := args.Get(0).(map[uuid.UUID]bool)
	return res, args.Error(1)
}
