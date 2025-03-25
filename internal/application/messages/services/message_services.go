package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
)

type (
	IConversation interface {
		GetConversationById(ctx context.Context, conversationId uuid.UUID) (result *common.ConversationResult, err error)
		CreateConversation(ctx context.Context, command *command.CreateConversationCommand) (result *command.CreateConversationResult, err error)
		GetManyConversation(ctx context.Context, query *query.GetManyConversationQuery) (result *query.GetManyConversationQueryResult, err error)
		DeleteConversationById(ctx context.Context, command *command.DeleteConversationCommand) error
	}
	IMessage interface {
		GetMessageById(ctx context.Context, messageId uuid.UUID) (result *common.MessageResult, err error)
		CreateMessage(ctx context.Context, command *command.CreateMessageCommand) error
		GetMessagesByConversationId(ctx context.Context, query *query.GetMessagesByConversationIdQuery) (result *query.GetMessagesByConversationIdResult, err error)
		DeleteMessageById(ctx context.Context, command *command.DeleteMessageCommand) error
	}
	IConversationDetail interface {
		GetConversationDetailById(ctx context.Context, userId uuid.UUID, conversationId uuid.UUID) (result *common.ConversationDetailResult, err error)
		CreateConversationDetail(ctx context.Context, command *command.CreateConversationDetailCommand) (result *command.CreateConversationDetailResult, err error)
		GetConversationDetailByIdList(ctx context.Context, query *query.GetConversationDetailQuery) (result *query.GetConversationDetailResult, err error)
		DeleteConversationDetailById(ctx context.Context, command *command.DeleteConversationDetailCommand) error
		UpdateOneStatusConversationDetail(ctx context.Context, command *command.UpdateOneStatusConversationDetailCommand) (err error)
	}
	IMessageMQ interface {
		HandleMessage(ctx context.Context, message *command.CreateMessageCommand) error
	}
)

var (
	localConversation       IConversation
	localMessage            IMessage
	localConversationDetail IConversationDetail
	localMessageMQ          IMessageMQ
)

func Conversation() IConversation {
	if localConversation == nil {
		panic("service_implement localConversation not found for interface IConversation")
	}
	return localConversation
}

func InitConversation(i IConversation) {
	localConversation = i
}

func Message() IMessage {
	if localMessage == nil {
		panic("service_implement localMessage not found for interface IMessage")
	}
	return localMessage
}

func InitMessage(i IMessage) {
	localMessage = i
}

func ConversationDetail() IConversationDetail {
	if localConversationDetail == nil {
		panic("service_implement localConversationDetail not found for interface IConversationDetail")
	}
	return localConversationDetail
}

func InitConversationDetail(i IConversationDetail) {
	localConversationDetail = i
}

func MessageMQ() IMessageMQ {
	if localMessageMQ == nil {
		panic("service_implement localMessageMQ not found for interface IMessageMQ")
	}
	return localMessageMQ
}

func InitMessageMQ(i IMessageMQ) {
	localMessageMQ = i
}
