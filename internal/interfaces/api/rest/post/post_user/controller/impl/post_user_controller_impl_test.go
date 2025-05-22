package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	PostQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostUserService struct {
	mock.Mock
}

func (m *MockPostUserService) CreatePost(ctx context.Context, cmd *command.CreatePostCommand) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

func (m *MockPostUserService) ApprovePost(ctx context.Context, command *command.ApprovePostCommand) error {
	return nil
}

func (m *MockPostUserService) RejectPost(ctx context.Context, command *command.RejectPostCommand) error {
	return nil
}

func (m *MockPostUserService) UpdatePost(ctx context.Context, command *command.UpdatePostCommand) (*command.UpdatePostCommandResult, error) {
	return nil, nil
}

func (m *MockPostUserService) DeletePost(ctx context.Context, command *command.DeletePostCommand) error {
	return nil
}

func (m *MockPostUserService) GetPost(ctx context.Context, query *PostQuery.GetOnePostQuery) (*PostQuery.GetOnePostQueryResult, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*PostQuery.GetOnePostQueryResult), args.Error(1)
}

func (m *MockPostUserService) GetManyPosts(ctx context.Context, query *PostQuery.GetManyPostQuery) (*PostQuery.GetManyPostQueryResult, error) {
	return nil, nil
}

func (m *MockPostUserService) GetTrendingPost(ctx context.Context, query *PostQuery.GetTrendingPostQuery) (*PostQuery.GetManyPostQueryResult, error) {
	return nil, nil
}

func (m *MockPostUserService) CheckPostOwner(ctx context.Context, query *PostQuery.CheckPostOwnerQuery) (*PostQuery.CheckPostOwnerQueryResult, error) {
	return nil, nil
}

func (m *MockPostUserService) ClearAllPostCaches(ctx context.Context) error {
	return nil
}

func TestCreatePost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Tạo mock service
	mockService := new(MockPostUserService)
	mockService.On("CreatePost", mock.Anything, mock.MatchedBy(func(cmd *command.CreatePostCommand) bool {
		// Kiểm tra các giá trị bên trong CreatePostCommand
		return cmd.UserId.String() != uuid.Nil.String() &&
			cmd.Content == "Hello world" &&
			cmd.Privacy == consts.PUBLIC
	})).Return(nil)

	// 2. Tạo controller với service mock
	postController := NewPostUserController(mockService)

	// 3. Tạo router và middleware
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Tạo userId dạng uuid.UUID
		userID := uuid.New()
		c.Set("userId", userID)

		// Tạo fake CreatePostRequest
		c.Set("validatedRequest", &request.CreatePostRequest{
			Content:  "Hello world",
			Privacy:  consts.PUBLIC,
			Location: "Vietnam",
		})
		c.Next()
	})
	router.Use(middlewares.ErrorHandlerMiddleware())

	// 4. Gắn handler
	router.POST("/posts", postController.CreatePost)

	// 5. Gửi request giả lập
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("content", "Hello world")
	_ = writer.WriteField("privacy", string(consts.PUBLIC))
	writer.Close()

	req := httptest.NewRequest("POST", "/posts", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer mock-token") // Thêm token giả lập

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 6. Kiểm tra kết quả
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Success")

	// 7. Xác nhận CreatePost đã được gọi với các tham số đúng
	mockService.AssertCalled(t, "CreatePost", mock.Anything, mock.MatchedBy(func(cmd *command.CreatePostCommand) bool {
		// Kiểm tra lại các giá trị bên trong
		return cmd.Content == "Hello world" && cmd.Privacy == consts.PUBLIC
	}))
}

func TestCreatePost_ValidationFail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Tạo mock service nhưng không cần gán kỳ vọng vì sẽ không được gọi
	mockService := new(MockPostUserService)
	mockService.On("CreatePost", mock.Anything, mock.Anything).Return(nil)

	// Tạo controller
	postController := NewPostUserController(mockService)

	// Tạo router
	router := gin.New()

	// Gắn middleware xử lý lỗi (nếu có)
	router.Use(middlewares.ErrorHandlerMiddleware())
	// Gắn route handler
	router.POST("/posts",
		helpers.ValidateFormBody(&request.CreatePostRequest{}, request.ValidateCreatePostRequest), postController.CreatePost)

	// 5. Tạo request giả lập
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("content", "1")
	writer.WriteField("privacy", string(consts.PUBLIC))
	writer.Close()

	req := httptest.NewRequest("POST", "/posts", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer mock-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 6. Kiểm tra kết quả trả về
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	errorData := response["error"].(map[string]interface{})
	assert.Equal(t, float64(50004), errorData["code"])
	assert.Equal(t, "Validate param failed", errorData["message"])
	assert.Contains(t, errorData["message_detail"], "Content: the length must be between 2 and 10000.")

	// 7. Xác nhận rằng CreatePost KHÔNG được gọi
	mockService.AssertNotCalled(t, "CreatePost", mock.Anything, mock.Anything)
}

func TestGetPostById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	postID := uuid.New()
	userID := uuid.New()

	// Mock kết quả post
	mockPost := &common.PostResultWithLiked{
		ID:           postID,
		UserId:       userID,
		Content:      "This is a test post",
		LikeCount:    10,
		CommentCount: 5,
		Status:       true,
		Location:     "Hanoi",
		IsLiked:      true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		User: &common.UserResult{
			ID:         userID,
			FamilyName: "testuser",
			AvatarUrl:  "http://example.com/avatar.jpg",
			Name:       "Test User",
		},
		// Bổ sung các field khác nếu cần test kỹ hơn
	}

	// Mock service
	mockService := new(MockPostUserService)
	mockService.
		On("GetPost", mock.Anything, mock.AnythingOfType("*query.GetOnePostQuery")).
		Return(&PostQuery.GetOnePostQueryResult{
			Post: mockPost,
		}, nil)

	// Controller
	postController := NewPostUserController(mockService)

	// Router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", userID)
		c.Next()
	})
	router.Use(middlewares.ErrorHandlerMiddleware())

	// Route gắn đúng
	router.GET("/posts/:post_id", postController.GetPostById)

	// Request
	req := httptest.NewRequest("GET", "/posts/"+postID.String(), nil)
	req.Header.Set("Authorization", "Bearer mock-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Kiểm tra kết quả
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "This is a test post")
	assert.Contains(t, w.Body.String(), postID.String())
	assert.Contains(t, w.Body.String(), "like_count")
	assert.Contains(t, w.Body.String(), "comment_count")

	mockService.AssertCalled(t, "GetPost", mock.Anything, mock.AnythingOfType("*query.GetOnePostQuery"))
}

func TestGetPostById_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()

	// Mock service không cần gọi, vì phải validate ngay từ router/controller

	postController := NewPostUserController(nil) // nil vì service không được gọi

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", userID)
		c.Next()
	})
	router.Use(middlewares.ErrorHandlerMiddleware())
	router.GET("/posts/:post_id", postController.GetPostById)

	// Truyền post_id không hợp lệ
	req := httptest.NewRequest("GET", "/posts/invalid-uuid-12345", nil)
	req.Header.Set("Authorization", "Bearer mock-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Kiểm tra trả về 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid UUID")
}
