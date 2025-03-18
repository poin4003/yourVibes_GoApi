package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/global"
	messageCommand "github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"go.uber.org/zap"
)

type sMessageMQ struct {
	conversationDetailRepo repositories.IConversationDetailRepository
}

func NewMessageMQImplement(
	conversationDetailRepo repositories.IConversationDetailRepository,
) *sMessageMQ {
	return &sMessageMQ{
		conversationDetailRepo: conversationDetailRepo,
	}
}

func (s *sMessageMQ) HandleMessage(
	ctx context.Context,
	command *messageCommand.CreateMessageCommand,
) error {
	userIds, err := s.conversationDetailRepo.GetListUserIdByConversationId(ctx, command.ConversationId)
	if err != nil {
		global.Logger.Error("GetListUserIdByConversationId", zap.Error(err))
	}

	if global.MessageSocketHub != nil {
		for _, userId := range userIds {
			if err = global.MessageSocketHub.SendMessage(userId.String(), command); err != nil {
				global.Logger.Error("SendNotification", zap.Error(err))
			}
		}
	}
	return nil
}
