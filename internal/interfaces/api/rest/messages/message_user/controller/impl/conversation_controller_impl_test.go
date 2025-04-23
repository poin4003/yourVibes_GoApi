package impl

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConversationService struct {
	mock.Mock
}

func (m *MockConversationService) CreateConversation(ctx context.Context, cmd *command.CreateConversationCommand) (*command.CreateConversationResult, error) {
	args := m.Called(ctx, cmd)
	if res, ok := args.Get(0).(*command.CreateConversationResult); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockConversationService) GetConversationById(ctx context.Context, conversationId uuid.UUID) (*common.ConversationResult, error) {
	args := m.Called(ctx, conversationId)
	if res, ok := args.Get(0).(*common.ConversationResult); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockConversationService) GetManyConversation(ctx context.Context, userId uuid.UUID, query *query.GetManyConversationQuery) (*query.GetManyConversationQueryResult, error) {
	// args := m.Called(ctx, userId, query)
	// if res, ok := args.Get(0).(*query.GetManyConversationQueryResult); ok {
	// 	return res, args.Error(1)
	// }
	// return nil, args.Error(1)
	return nil, nil
}

func (m *MockConversationService) DeleteConversationById(ctx context.Context, cmd *command.DeleteConversationCommand) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

func (m *MockConversationService) UpdateConversationById(ctx context.Context, cmd *command.UpdateConversationCommand) (*command.UpdateConversationCommandResult, error) {
	args := m.Called(ctx, cmd)
	if res, ok := args.Get(0).(*command.UpdateConversationCommandResult); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestCreateConversation_Success(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Test Chat")
	writer.WriteField("user_ids", uuid.New().String())
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/conversations", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", uuid.New())
	c.Set("validatedRequest", &request.CreateConversationRequest{
		Name:    "Test Chat",
		UserIds: []string{uuid.New().String()},
	})

	expected := &command.CreateConversationResult{
		Conversation: &common.ConversationResult{
			ID:   uuid.New(),
			Name: "Test Chat",
		},
	}
	mockService.On("CreateConversation", mock.Anything, mock.AnythingOfType("*command.CreateConversationCommand")).Return(expected, nil)

	controller.CreateConversation(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateConversation_Fail(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "")
	writer.WriteField("user_ids", "")
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/conversations", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", uuid.New())
	c.Set("validatedRequest", &request.CreateConversationRequest{
		Name:    "",
		UserIds: []string{},
	})

	mockService.On("CreateConversation", mock.Anything, mock.AnythingOfType("*command.CreateConversationCommand")).Return(nil, errors.New("cannot create"))

	controller.CreateConversation(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Missing validatedRequest request")
}

func TestGetConversationById_Success(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/conversations/{id}", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", uuid.New())

	mockService.On("GetConversationById", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(&common.ConversationResult{
		ID:   uuid.New(),
		Name: "Test Chat",
	}, nil)

	controller.GetConversationById(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetConversationById_Fail_NotFound(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	id := uuid.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/conversations/"+id.String(), nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "conversationId", Value: id.String()}}
	mockService.On("GetConversationById", mock.Anything, id).Return(nil, errors.New("not found"))

	controller.GetConversationById(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteConversationById_Success(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	id := uuid.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/conversations/"+id.String(), nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "conversationId", Value: id.String()}}
	mockService.On("DeleteConversationById", mock.Anything, mock.AnythingOfType("*command.DeleteConversationCommand")).Return(nil)

	controller.DeleteConversationById(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteConversationById_Fail(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	id := uuid.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/conversations/"+id.String(), nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "conversationId", Value: id.String()}}
	mockService.On("DeleteConversationById", mock.Anything, mock.AnythingOfType("*command.DeleteConversationCommand")).Return(errors.New("delete failed"))

	controller.DeleteConversationById(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateConversationById_Success(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	id := uuid.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/conversations/"+id.String(), nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "conversationId", Value: id.String()}}
	name := "Updated Name"
	c.Set("validatedRequest", &request.UpdateConversationRequest{Name: &name})

	expected := &command.UpdateConversationCommandResult{
		Conversation: &common.ConversationResult{
			ID:   id,
			Name: "Updated Name",
		},
	}
	mockService.On("UpdateConversationById", mock.Anything, mock.AnythingOfType("*command.UpdateConversationCommand")).Return(expected, nil)

	controller.UpdateConversation(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateConversationById_Fail(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	id := uuid.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/conversations/"+id.String(), nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "conversationId", Value: id.String()}}
	name := "Error Update"
	c.Set("validatedRequest", &request.UpdateConversationRequest{Name: &name})

	mockService.On("UpdateConversationById", mock.Anything, mock.AnythingOfType("*command.UpdateConversationCommand")).Return(nil, errors.New("update failed"))

	controller.UpdateConversation(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetManyConversation_Success(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	userID := uuid.New()
	queryData := &query.GetManyConversationQuery{Limit: 10, Page: 1}
	expected := &query.GetManyConversationQueryResult{
		Conversation: []*common.ConversationResult{
			{ID: uuid.New(), Name: "Chat 1"},
			{ID: uuid.New(), Name: "Chat 2"},
		},
		PagingResponse: &response.PagingResponse{
			Total: 2,
			Limit: 10,
			Page:  1,
		},
	}

	mockService.On("GetManyConversation", mock.Anything, userID, queryData).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/conversations", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID)
	c.Set("query", queryData)

	controller.GetConversation(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetManyConversation_Fail(t *testing.T) {
	mockService := new(MockConversationService)
	controller := NewConversationController(mockService)

	userID := uuid.New()
	queryData := &query.GetManyConversationQuery{Limit: 10, Page: 1}
	mockService.On("GetManyConversation", mock.Anything, userID, queryData).Return(nil, errors.New("failed to fetch"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/conversations", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID)
	c.Set("query", queryData)

	controller.GetConversation(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
