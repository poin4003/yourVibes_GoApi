package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"gorm.io/gorm"
)

type CreateMessageCommand struct {
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
	MediaUrls      []string
}

type CreateMessageCommandHandler struct {
	db *gorm.DB
}

func NewCreateMessageCommandHandler(db *gorm.DB) *CreateMessageCommandHandler {
	return &CreateMessageCommandHandler{db: db}
}

func (h *CreateMessageCommandHandler) Handle(ctx context.Context, cmd CreateMessageCommand) (*models.Message, error) {
	// 1. Kiểm tra tính hợp lệ của Command.
	if cmd.SenderID == uuid.Nil {
		return nil, errors.New("sender ID is required")
	}
	if cmd.ConversationID == uuid.Nil {
		return nil, errors.New("conversation ID is required")
	}

	// 2. Kiểm tra xem cuộc trò chuyện có tồn tại hay không.
	var conversation models.Conversation
	if err := h.db.WithContext(ctx).First(&conversation, cmd.ConversationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("conversation with ID %s not found: %w", cmd.ConversationID, err)
		}
		return nil, fmt.Errorf("failed to find conversation: %w", err)
	}

	// 3. Kiểm tra xem người gửi có phải là thành viên của cuộc trò chuyện hay không.
	var member models.ConversationMember
	if err := h.db.WithContext(ctx).Where("conversation_id = ? AND user_id = ?", cmd.ConversationID, cmd.SenderID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sender is not a member of the conversation: %w", err)
		}
		return nil, fmt.Errorf("failed to check membership: %w", err)
	}

	// 4. Tạo tin nhắn.
	newMessage := &models.Message{
		ConversationID: cmd.ConversationID,
		SenderID:       cmd.SenderID,
		Content:        cmd.Content,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 5. Tạo và đính kèm Media.
	if len(cmd.MediaUrls) > 0 {
		for _, mediaURL := range cmd.MediaUrls {
			newMessage.Media = append(newMessage.Media, models.MessageMedia{
				MediaUrl:  mediaURL,
				Status:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	// 6. Lưu tin nhắn vào Cơ sở dữ liệu.
	if err := h.db.WithContext(ctx).Create(newMessage).Error; err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return newMessage, nil
}
