package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/global"
	messageCommand "github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/socket_hub"
	"go.uber.org/zap"
)

type sMessageMQ struct {
	conversationDetailRepo repositories.IConversationDetailRepository
	messageSocketHub       *socket_hub.MessageSocketHub
}

func NewMessageMQImplement(
	conversationDetailRepo repositories.IConversationDetailRepository,
	messageSocketHub *socket_hub.MessageSocketHub,
) *sMessageMQ {
	return &sMessageMQ{
		conversationDetailRepo: conversationDetailRepo,
		messageSocketHub:       messageSocketHub,
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

	for _, userId := range userIds {
		if err = s.messageSocketHub.SendMessage(userId.String(), command); err != nil {
			global.Logger.Error("SendNotification", zap.Error(err))
		}
	}

	return nil
}
