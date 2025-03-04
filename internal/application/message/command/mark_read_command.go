package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"gorm.io/gorm"
)
 
type MarkReadMessageCommand struct {
	MessageID uuid.UUID  
	UserID uuid.UUID  
}
type MarkReadMessageCommandHandler struct {
	db *gorm.DB 
}

func NewMarkReadMessageCommandHandler(db *gorm.DB) *MarkReadMessageCommandHandler {
	return &MarkReadMessageCommandHandler{db: db}
}

func (h *MarkReadMessageCommandHandler) Handle(ctx context.Context, cmd MarkReadMessageCommand) (*models.Message, error) {
	// 1. Kiểm tra tính hợp lệ của Command.
	if cmd.MessageID == uuid.Nil {
		return nil, errors.New("message ID is required") 
	}
	if cmd.UserID == uuid.Nil {
		return nil, errors.New("user ID is required") 
	}

	// 2. Kiểm tra xem tin nhắn có tồn tại hay không.
	var message models.Message
	if err := h.db.WithContext(ctx).Preload("Conversation").First(&message, cmd.MessageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("message with ID %s not found: %w", cmd.MessageID, err) 
		}
		return nil, fmt.Errorf("failed to find message: %w", err) 
	}

	//3. Kiểm tra người dùng có là thành viên của cuộc trò chuyện hay không.
	var member models.ConversationMember
	if err := h.db.WithContext(ctx).Where("conversation_id = ? AND user_id = ?", message.ConversationID, cmd.UserID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user is not a member of the conversation: %w", err) 
		}
		return nil, fmt.Errorf("failed to check membership: %w", err) 
	}

	// 4. Cập nhật trạng thái tin nhắn thành false (đã đọc).
	if err := h.db.WithContext(ctx).Model(&message).Update("status", false).Error; err != nil {
		return nil, fmt.Errorf("failed to mark message as read: %w", err)
	}

	return &message, nil 
}
