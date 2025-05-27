package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/stretchr/testify/mock"
)

type MockFriendRepository struct {
	mock.Mock
}

func (m *MockFriendRepository) CreateOne(ctx context.Context, entity *entities.Friend) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockFriendRepository) DeleteOne(ctx context.Context, entity *entities.Friend) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockFriendRepository) GetFriends(ctx context.Context, query *query.FriendQuery) ([]*entities.User, *response.PagingResponse, error) {
	args := m.Called(ctx, query)
	users, _ := args.Get(0).([]*entities.User)
	paging, _ := args.Get(1).(*response.PagingResponse)
	return users, paging, args.Error(2)
}

func (m *MockFriendRepository) GetFriendIds(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, userId)
	ids, _ := args.Get(0).([]uuid.UUID)
	return ids, args.Error(1)
}

func (m *MockFriendRepository) CheckFriendExist(ctx context.Context, entity *entities.Friend) (bool, error) {
	args := m.Called(ctx, entity)
	return args.Bool(0), args.Error(1)
}

func (m *MockFriendRepository) GetFriendSuggestions(ctx context.Context, query *query.FriendQuery) ([]*entities.UserWithSendFriendRequest, *response.PagingResponse, error) {
	args := m.Called(ctx, query)
	users, _ := args.Get(0).([]*entities.UserWithSendFriendRequest)
	paging, _ := args.Get(1).(*response.PagingResponse)
	return users, paging, args.Error(2)
}

func (m *MockFriendRepository) GetFriendByBirthday(ctx context.Context, query *query.FriendQuery) ([]*entities.UserWithBirthday, *response.PagingResponse, error) {
	args := m.Called(ctx, query)
	users, _ := args.Get(0).([]*entities.UserWithBirthday)
	paging, _ := args.Get(1).(*response.PagingResponse)
	return users, paging, args.Error(2)
}

func (m *MockFriendRepository) GetNonFriends(ctx context.Context, query *query.FriendQuery) ([]*entities.UserWithSendFriendRequest, *response.PagingResponse, error) {
	args := m.Called(ctx, query)
	users, _ := args.Get(0).([]*entities.UserWithSendFriendRequest)
	paging, _ := args.Get(1).(*response.PagingResponse)
	return users, paging, args.Error(2)
}
