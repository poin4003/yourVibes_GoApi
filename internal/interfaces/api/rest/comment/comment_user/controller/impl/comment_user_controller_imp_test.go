package impl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	CommentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	CommentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommentUserService struct {
	mock.Mock
}

func (m *MockCommentUserService) CreateComment(ctx context.Context, command *CommentCommand.CreateCommentCommand) (*CommentCommand.CreateCommentResult, error) {
	args := m.Called(ctx, command)
	var result *CommentCommand.CreateCommentResult
	if args.Get(0) != nil {
		result = args.Get(0).(*CommentCommand.CreateCommentResult)
	}
	return result, args.Error(1)
}

func (m *MockCommentUserService) UpdateComment(ctx context.Context, command *CommentCommand.UpdateCommentCommand) (*CommentCommand.UpdateCommentResult, error) {
	args := m.Called(ctx, command)
	var result *CommentCommand.UpdateCommentResult
	if args.Get(0) != nil {
		result = args.Get(0).(*CommentCommand.UpdateCommentResult)
	}
	return result, args.Error(1)
}

func (m *MockCommentUserService) DeleteComment(ctx context.Context, command *CommentCommand.DeleteCommentCommand) error {
	args := m.Called(ctx, command)
	return args.Error(0)
}

func (m *MockCommentUserService) GetManyComments(ctx context.Context, query *CommentQuery.GetManyCommentQuery) (*CommentQuery.GetManyCommentsResult, error) {
	args := m.Called(ctx, query)
	var result *CommentQuery.GetManyCommentsResult
	if args.Get(0) != nil {
		result = args.Get(0).(*CommentQuery.GetManyCommentsResult)
	}
	return result, args.Error(1)
}

func (m *MockCommentUserService) ClearAllCommentCaches(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestCreateComment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	postID := uuid.New()

	mockService := new(MockCommentUserService)
	commentController := NewCommentUserController(mockService)

	// Tạo request giả lập
	requestData := &request.CreateCommentRequest{
		Content: "Test comment",
		PostId:  postID,
	}

	// Convert sang command
	expectedCommand, _ := requestData.ToCreateCommentCommand(userID)

	// Tạo comment entity
	expectedEntity, err := entities.NewComment(
		uuid.New(), // CommentId
		expectedCommand.PostId,
		expectedCommand.UserId,
		expectedCommand.ParentId,
		expectedCommand.Content,
	)
	assert.NoError(t, err)

	// Trả về từ mock service
	mockService.On("CreateComment", mock.Anything, expectedCommand).Return(&CommentCommand.CreateCommentResult{
		Comment: &common.CommentResult{
			ID:              expectedEntity.ID,
			PostId:          expectedEntity.PostId,
			UserId:          expectedEntity.UserId,
			ParentId:        expectedEntity.ParentId,
			Content:         expectedEntity.Content,
			LikeCount:       expectedEntity.LikeCount,
			RepCommentCount: expectedEntity.RepCommentCount,
			CreatedAt:       expectedEntity.CreatedAt,
			UpdatedAt:       expectedEntity.UpdatedAt,
			Status:          expectedEntity.Status,
			User: &common.UserResult{
				ID:         expectedEntity.UserId,
				Name:       "Test User",
				FamilyName: "User",
				AvatarUrl:  "http://example.com/avatar.jpg",
			},
		}}, nil)

	router := gin.Default()

	// Middleware giả lập set context
	router.POST("/comments/", func(ctx *gin.Context) {
		ctx.Set("userId", userID)
		ctx.Set("validatedRequest", requestData)
		commentController.CreateComment(ctx)
	})

	req := httptest.NewRequest("POST", "/comments/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Kiểm tra dữ liệu trả về
	expectedResponse := response.ToCommentDto(&common.CommentResult{
		ID:              expectedEntity.ID,
		PostId:          expectedEntity.PostId,
		UserId:          expectedEntity.UserId,
		ParentId:        expectedEntity.ParentId,
		Content:         expectedEntity.Content,
		LikeCount:       expectedEntity.LikeCount,
		RepCommentCount: expectedEntity.RepCommentCount,
		CreatedAt:       expectedEntity.CreatedAt,
		UpdatedAt:       expectedEntity.UpdatedAt,
		Status:          expectedEntity.Status,
		User: &common.UserResult{
			ID:         expectedEntity.UserId,
			Name:       "Test User",
			FamilyName: "User",
			AvatarUrl:  "http://example.com/avatar.jpg",
		},
	})

	// Check nội dung trả về
	body := w.Body.String()
	assert.Contains(t, body, expectedResponse.ID.String())
	assert.Contains(t, body, expectedResponse.Content)
}

func TestCreateComment_MissingContext_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	mockService := new(MockCommentUserService)
	commentController := NewCommentUserController(mockService)

	router := gin.Default()

	router.POST("/comments/", func(ctx *gin.Context) {
		ctx.Set("userId", userID)
		// KHÔNG set "validatedRequest"
		commentController.CreateComment(ctx)
	})

	req := httptest.NewRequest("POST", "/comments/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing validatedRequest")
}

func TestCreateComment_InvalidPostId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockCommentUserService)
	mockService.On("CreateComment", mock.Anything, mock.Anything).Return(nil, nil)
	commentController := NewCommentUserController(mockService)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", uuid.New())
		c.Next()
	})
	router.Use(middlewares.ErrorHandlerMiddleware())

	router.POST("/comments/",
		helpers.ValidateJsonBody(&request.CreateCommentRequest{}, request.ValidateCreateCommentRequest), commentController.CreateComment)

	body := `{"content": "Test comment", "postId": "invalid-uuid"}`
	req := httptest.NewRequest("POST", "/comments/", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer mock-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid postId")

	// Đảm bảo không gọi mock service
	mockService.AssertNotCalled(t, "CreateComment")
}

func TestGetComment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 2. Chuẩn bị dữ liệu input
	userID := uuid.New()
	postID := uuid.New()

	mockComment := &common.CommentResultWithLiked{
		ID:              uuid.New(),
		PostId:          postID,
		UserId:          userID,
		User:            &common.UserResult{ID: userID, Name: "Test User", FamilyName: "User", AvatarUrl: "http://example.com/avatar.jpg"},
		ParentId:        nil,
		Content:         "Test comment",
		LikeCount:       0,
		RepCommentCount: 0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          true,
		IsLiked:         true,
	}

	mockService := new(MockCommentUserService)
	mockService.On("GetManyComments", mock.Anything, mock.MatchedBy(func(query *CommentQuery.GetManyCommentQuery) bool {
		return query.PostId.String() == postID.String() &&
			query.Limit == 10 &&
			query.Page == 1
	})).Return(&CommentQuery.GetManyCommentsResult{
		Comments: []*common.CommentResultWithLiked{mockComment},
		PagingResponse: &pkgResponse.PagingResponse{
			Total: 1,
			Limit: 10,
			Page:  1,
		},
	}, nil)
	commentController := NewCommentUserController(mockService)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", userID)
		c.Set("validateQuery", &CommentQuery.GetManyCommentQuery{
			PostId: postID,
			Limit:  10,
			Page:   1,
		})
		c.Next()
	})
	router.Use(middlewares.ErrorHandlerMiddleware())
	router.GET("/comments/", helpers.ValidateQuery(&query.CommentQueryObject{}, query.ValidateCommentQueryObject), commentController.GetComment)
	req := httptest.NewRequest("GET", "/comments/", nil)
	req.Header.Set("Authorization", "Bearer mock-token") // Thêm token giả lập
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Kiểm tra kết quả
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test comment") // Hoặc kiểm tra các trường khác nếu cần thiết
	// Xác nhận GetManyComments được gọi
	mockService.AssertCalled(t, "GetManyComments", mock.Anything, mock.AnythingOfType("*query.GetManyCommentQuery"))

}
