package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/message/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/message/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	IMessageRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Message, error)
		CreateOne(ctx context.Context, entity *entities.Message) (*entities.Message, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.MessageUpdate) (*entities.Message, error)
		UpdateMany(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteOne(ctx context.Context, id uuid.UUID) (*entities.Message, error)
		DeleteMany(ctx context.Context, condition map[string]interface{}) (int64, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Message, error)
		GetMany(ctx context.Context, query *query.GetManyMessageQuery) ([]*entities.Message, *response.PagingResponse, error)
		GetConversationMessages(ctx context.Context, conversationId uuid.UUID, query *query.GetManyMessageQuery) ([]*entities.Message, *response.PagingResponse, error)
	}

	IConversationRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Conversation, error)
		CreateOne(ctx context.Context, entity *entities.Conversation) (*entities.Conversation, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.ConversationUpdate) (*entities.Conversation, error)
		UpdateMany(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteOne(ctx context.Context, id uuid.UUID) (*entities.Conversation, error)
		DeleteMany(ctx context.Context, condition map[string]interface{}) (int64, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Conversation, error)
		GetMany(ctx context.Context, query *query.GetManyConversationQuery) ([]*entities.Conversation, *response.PagingResponse, error)
		CheckConversationExist(ctx context.Context, conversationId uuid.UUID) (bool, error)
	}

	IConversationUserRepository interface {
		CreateOne(ctx context.Context, entity *entities.ConversationUser) (*entities.ConversationUser, error)
		CreateMany(ctx context.Context, entity []*entities.ConversationUser) ([]*entities.ConversationUser, error)
		DeleteOne(ctx context.Context, conversationId uuid.UUID, userId uuid.UUID) (*entities.ConversationUser, error)
		DeleteMany(ctx context.Context, conversationId uuid.UUID) (int64, error)
		GetMany(ctx context.Context, query *query.GetManyConversationUserQuery) ([]*entities.ConversationUser, *response.PagingResponse, error)
		CheckExist(ctx context.Context, conversationId uuid.UUID, userId uuid.UUID) (bool, error)
		GetByConversationId(ctx context.Context, conversationId uuid.UUID) ([]*entities.ConversationUser, error)
	}
)

var (
	localMessage          IMessageRepository
	localConversation     IConversationRepository
	localConversationUser IConversationUserRepository
)

func Message() IMessageRepository {
	if localMessage == nil {
		panic("repository_implement localMessage not found for interface IMessage")
	}

	return localMessage
}

func Conversation() IConversationRepository {
	if localConversation == nil {
		panic("repository_implement localConversation not found for interface IConversation")
	}

	return localConversation
}

func ConversationUser() IConversationUserRepository {
	if localConversationUser == nil {
		panic("repository_implement localConversationUser not found for interface IConversationUser")
	}
	return localConversationUser
}

func InitMessageRepository(i IMessageRepository) {
	localMessage = i
}

func InitConversationRepository(i IConversationRepository) {
	localConversation = i
}

func InitConversationUserRepository(i IConversationUserRepository) {
	localConversationUser = i
}
