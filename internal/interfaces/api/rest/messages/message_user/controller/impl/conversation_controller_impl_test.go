package impl

// import (
// 	"context"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
// 	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
// 	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // Mock service implements IConversation
// type MockConversationService struct {
// 	mock.Mock
// }

// func (m *MockConversationService) GetConversationById(ctx context.Context, conversationId uuid.UUID) (*common.ConversationResult, error) {
// 	args := m.Called(ctx, conversationId)
// 	if args.Get(0) != nil {
// 		return args.Get(0).(*common.ConversationResult), args.Error(1)
// 	}
// 	return nil, args.Error(1)
// }

// func (m *MockConversationService) CreateConversation(ctx context.Context, cmd *command.CreateConversationCommand) (*command.CreateConversationResult, error) {
// 	args := m.Called(ctx, cmd)
// 	if args.Get(0) != nil {
// 		return args.Get(0).(*command.CreateConversationResult), args.Error(1)
// 	}
// 	return nil, args.Error(1)
// }

// // func (m *MockConversationService) GetManyConversation(ctx context.Context, userId uuid.UUID, query *query.GetManyConversationQuery) (*query.GetManyConversationQueryResult, error) {
// // 	args := m.Called(ctx, userId, query)
// // 	if args.Get(0) != nil {
// // 		return args.Get(0).(*query.GetManyConversationQueryResult), args.Error(1)
// // 	}
// // 	return nil, args.Error(1)
// // }

// func (m *MockConversationService) DeleteConversationById(ctx context.Context, cmd *command.DeleteConversationCommand) error {
// 	args := m.Called(ctx, cmd)
// 	return args.Error(0)
// }

// func (m *MockConversationService) UpdateConversationById(ctx context.Context, cmd *command.UpdateConversationCommand) (*command.UpdateConversationCommandResult, error) {
// 	args := m.Called(ctx, cmd)
// 	if args.Get(0) != nil {
// 		return args.Get(0).(*command.UpdateConversationCommandResult), args.Error(1)
// 	}
// 	return nil, args.Error(1)
// }

// func TestGetConversationById(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	mockService := new(MockConversationService)
// 	controller := NewConversationController(mockService)

// 	t.Run("success", func(t *testing.T) {
// 		// Setup
// 		w := httptest.NewRecorder()
// 		ctx, _ := gin.CreateTestContext(w)
// 		conversationId := uuid.New()
// 		ctx.Set("conversationId", conversationId)

// 		expectedResult := &common.ConversationResult{
// 			ID:             conversationId,
// 			Name:           "Test Conversation",
// 			Image:          "",
// 			Avatar:         "",
// 			UserID:         nil,
// 			FamilyName:     "",
// 			CreatedAt:      time.Now(),
// 			UpdatedAt:      time.Now(),
// 			LastMess:       nil,
// 			LastMessStatus: false,
// 		}

// 		mockService.On("GetConversationById", mock.Anything, conversationId).Return(expectedResult, nil)

// 		// Act
// 		controller.GetConversationById(ctx)

// 		// Assert
// 		assert.Equal(t, 200, w.Code)
// 		mockService.AssertExpectations(t)
// 	})
// 	t.Run("fail - invalid conversation ID", func(t *testing.T) {
// 		// Setup
// 		w := httptest.NewRecorder()
// 		ctx, _ := gin.CreateTestContext(w)
// 		ctx.Params = gin.Params{gin.Param{Key: "conversationId", Value: "invalid-id"}}

// 		// Act
// 		controller.GetConversationById(ctx)

// 		// Assert
// 		assert.Equal(t, 400, w.Code)

// 	})
// 	t.Run("fail - conversation not found", func(t *testing.T) {
// 		// Setup
// 		w := httptest.NewRecorder()
// 		ctx, _ := gin.CreateTestContext(w)
// 		conversationId := uuid.New()
// 		ctx.Set("conversationId", conversationId)

// 		mockService.On("GetConversationById", mock.Anything, conversationId).Return(nil, assert.AnError)

// 		// Act
// 		controller.GetConversationById(ctx)

// 		// Assert
// 		assert.Equal(t, 404, w.Code)
// 		mockService.AssertExpectations(t)
// 	})

// 	t.Run("fail - service error", func(t *testing.T) {
// 		// Given
// 		conversationID := uuid.New()

// 		mockService.On("GetConversationById", mock.Anything, conversationID).Return(nil, assert.AnError)

// 		// Setup
// 		w := httptest.NewRecorder()
// 		ctx, _ := gin.CreateTestContext(w)
// 		ctx.Params = gin.Params{gin.Param{Key: "conversationId", Value: conversationID.String()}}

// 		// When
// 		controller.GetConversationById(ctx)

// 		// Then
// 		assert.Equal(t, 500, w.Code)
// 		mockService.AssertExpectations(t)
// 	})
// }

// func TestCreateConversation(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	mockService := new(MockConversationService)
// 	controller := NewConversationController(mockService)

// 	userId := uuid.New()
// 	userIds := []string{userId.String()}

// 	req := &request.CreateConversationRequest{
// 		Name:    "Test Conversation",
// 		UserIds: userIds,
// 	}

// 	cmd := &command.CreateConversationCommand{
// 		Name:    req.Name,
// 		Image:   "",
// 		UserIds: []uuid.UUID{userId, userId}, // ví dụ thêm chính mình và 1 người khác
// 	}

// 	expectedResult := &command.CreateConversationResult{
// 		Conversation: &common.ConversationResult{
// 			ID:             uuid.New(),
// 			Name:           "Test Conversation",
// 			Image:          "",
// 			Avatar:         "",
// 			UserID:         &userId,
// 			FamilyName:     "",
// 			CreatedAt:      time.Now(),
// 			UpdatedAt:      time.Now(),
// 			LastMess:       nil,
// 			LastMessStatus: false,
// 		},
// 	}

// 	t.Run("success", func(t *testing.T) {
// 		// Setup
// 		w := httptest.NewRecorder()
// 		ctx, _ := gin.CreateTestContext(w)
// 		ctx.Set("validatedRequest", req)
// 		ctx.Set("userId", userId)

// 		mockService.On("CreateConversation", mock.Anything, cmd).Return(expectedResult, nil)

// 		// Act
// 		controller.CreateConversation(ctx)

// 		// Assert
// 		assert.Equal(t, 200, w.Code)
// 		mockService.AssertExpectations(t)
// 	})

// 	t.Run("fail - service error", func(t *testing.T) {
// 		// Setup
// 		w := httptest.NewRecorder()
// 		ctx, _ := gin.CreateTestContext(w)
// 		ctx.Set("validatedRequest", req)
// 		ctx.Set("userId", userId)

// 		// Mocking service error
// 		mockService.On("CreateConversation", mock.Anything, cmd).Return(nil, assert.AnError)

// 		// Act
// 		controller.CreateConversation(ctx)

// 		// Assert
// 		// Check if the response code is 500 due to error in service
// 		assert.Equal(t, 500, w.Code)
// 		mockService.AssertExpectations(t)
// 	})
// }

// // func TestGetManyConversation(t *testing.T) {
// // 	gin.SetMode(gin.TestMode)

// // 	mockService := new(MockConversationService)
// // 	controller := NewConversationController(mockService)

// // 	userId := uuid.New()

// // 	queryRequest := &query.GetManyConversationQuery{
// // 		Page:  1,
// // 		Limit: 10,
// // 	}

// // 	ConversaionResult := &common.ConversationResult{
// // 		ID:             uuid.New(),
// // 		Name:           "Test Conversation",
// // 		Image:          "",
// // 		Avatar:         "",
// // 		UserID:         nil,
// // 		FamilyName:     "",
// // 		CreatedAt:      time.Now(),
// // 		UpdatedAt:      time.Now(),
// // 		LastMess:       nil,
// // 		LastMessStatus: false,
// // 	}

// // 	PagingResponse := &response.PagingResponse{
// // 		Page:  1,
// // 		Limit: 10,
// // 		Total: 1,
// // 	}

// // 	expectedResult := &query.GetManyConversationQueryResult{
// // 		Conversation:   []*common.ConversationResult{ConversaionResult},
// // 		PagingResponse: PagingResponse,
// // 	}

// // 	t.Run("success", func(t *testing.T) {
// // 		w := httptest.NewRecorder()
// // 		ctx, _ := gin.CreateTestContext(w)
// // 		ctx.Set("userId", userId)
// // 		ctx.Set("validatedQuery", queryRequest)

// // 		mockService.On("GetManyConversation", mock.Anything, userId, queryRequest).Return(expectedResult, nil)

// // 		controller.GetConversation(ctx)

// // 		assert.Equal(t, 200, w.Code)
// // 		mockService.AssertExpectations(t)
// // 	})

// // 	t.Run("fail - service error", func(t *testing.T) {
// // 		w := httptest.NewRecorder()
// // 		ctx, _ := gin.CreateTestContext(w)
// // 		ctx.Set("userId", userId)
// // 		ctx.Set("validatedQuery", queryRequest) // Sửa "query" thành "validatedQuery" để khớp với case success

// // 		mockService.On("GetManyConversation", mock.Anything, userId, queryRequest).Return(nil, assert.AnError)

// // 		controller.GetConversation(ctx)

// // 		assert.Equal(t, 500, w.Code)
// // 		mockService.AssertExpectations(t)
// // 	})
// // }
