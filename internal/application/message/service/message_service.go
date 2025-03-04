package implementations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/message"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/message/services"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"gorm.io/gorm"
)

// messageRepository chứa các repository liên quan đến tin nhắn.
type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository tạo một instance mới của messageRepository.
func NewMessageRepository(db *gorm.DB) *messageRepository {
	return &messageRepository{db: db}
}

// Implementations for message.MessageRepository
type messageImplementation struct {
	db *gorm.DB
}

// NewMessageImplementation creates a new MessageImplementation.
func NewMessageImplementation(db *gorm.DB) message.MessageRepository {
	return &messageImplementation{db: db}
}

func (r *messageImplementation) Create(ctx context.Context, msg *models.Message) (*models.Message, error) {
	if err := r.db.WithContext(ctx).Create(msg).Error; err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	return msg, nil
}

func (r *messageImplementation) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	var msg models.Message
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&msg).Error; err != nil {
		return nil, fmt.Errorf("failed to get message by ID: %w", err)
	}
	return &msg, nil
}

func (r *messageImplementation) Delete(ctx context.Context, msg *models.Message) error {
	if err := r.db.WithContext(ctx).Delete(msg).Error; err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

func (r *messageImplementation) DeleteByConversationID(ctx context.Context, conversationID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Delete(&models.Message{}).Error; err != nil {
		return fmt.Errorf("failed to delete messages by conversation ID: %w", err)
	}
	return nil
}

// Implementations for message.ConversationRepository
type conversationImplementation struct {
	db *gorm.DB
}

// NewConversationImplementation creates a new ConversationImplementation.
func NewConversationImplementation(db *gorm.DB) message.ConversationRepository {
	return &conversationImplementation{db: db}
}

func (r *conversationImplementation) Create(ctx context.Context, conversation *models.Conversation) (*models.Conversation, error) {
	if err := r.db.WithContext(ctx).Create(conversation).Error; err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	return conversation, nil
}

func (r *conversationImplementation) GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&conversation).Error; err != nil {
		return nil, fmt.Errorf("failed to get conversation by ID: %w", err)
	}
	return &conversation, nil
}

func (r *conversationImplementation) Delete(ctx context.Context, conversation *models.Conversation) error {
	if err := r.db.WithContext(ctx).Delete(conversation).Error; err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	return nil
}

// NewMessageService creates a new MessageService with dependencies.
func NewMessageService(db *gorm.DB) services.MessageService {

	// Repositories
	messageRepo := NewMessageImplementation(db)
	conversationRepo := NewConversationImplementation(db)
	messageUserRepo := serviceImplementations.NewMessageUserImplementation(db)
	messageLikeRepo := serviceImplementations.NewMessageLikeImplementation(db)
	// Command Handlers
	createMessageCommandHandler := command.NewCreateMessageCommandHandler(messageRepo, conversationRepo)

	// Query Handlers
	getMessagesQueryHandler := query.NewGetMessagesQueryHandler(db)
	getConversationsQueryHandler := query.NewGetConversationsQueryHandler(db)

	return services.NewMessageService(
		createMessageCommandHandler,
		getMessagesQueryHandler,
		getConversationsQueryHandler,
		messageUserRepo,
		messageLikeRepo,
	)
}
