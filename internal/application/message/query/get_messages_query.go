package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"gorm.io/gorm"
)

type GetMessagesQuery struct {
	ConversationID uuid.UUID 
	UserID         uuid.UUID 
	Limit          int       
	Offset         int       
}


type GetMessagesQueryResult struct {
	Messages []models.Message 
}


type GetMessagesQueryHandler struct {
	db *gorm.DB 
}


func NewGetMessagesQueryHandler(db *gorm.DB) *GetMessagesQueryHandler {
	return &GetMessagesQueryHandler{db: db}
}


func (h *GetMessagesQueryHandler) Handle(ctx context.Context, query GetMessagesQuery) (*GetMessagesQueryResult, error) {
	// 1. Kiểm tra tính hợp lệ của Query.
	if query.ConversationID == uuid.Nil {
		return nil, errors.New("conversation ID is required") 
	}
	if query.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}

	// 2. Kiểm tra cuộc trò chuyện có tồn tại không.
	var conversation models.Conversation
	if err := h.db.WithContext(ctx).First(&conversation, query.ConversationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("conversation with ID %s not found: %w", query.ConversationID, err) 
		}
		return nil, fmt.Errorf("failed to find conversation: %w", err)
	}

	//3. Kiểm tra người dùng có là thành viên của cuộc trò chuyện hay không.
	var member models.ConversationMember
	if err := h.db.WithContext(ctx).Where("conversation_id = ? AND user_id = ?", query.ConversationID, query.UserID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user is not a member of the conversation: %w", err) 
		}
		return nil, fmt.Errorf("failed to check membership: %w", err) 
	}

	// 4. Lấy danh sách tin nhắn.
	var messages []models.Message
	if err := h.db.WithContext(ctx).
		Where("conversation_id = ?", query.ConversationID).
		Order("created_at DESC").
		Limit(query.Limit).      
		Offset(query.Offset).   
		Preload("Sender").       
		Preload("Media").         
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err) 
	}

	// 5. Tạo kết quả trả về.
	result := &GetMessagesQueryResult{
		Messages: messages,
	}

	return result, nil 
}
