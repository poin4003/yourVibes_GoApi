package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/media"
	"sync"

	"github.com/google/uuid"
	conversationCommand "github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/mapper"
	conversationQuery "github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	conversationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sConversation struct {
	conversationRepo repositories.IConversationRepository
	userCache        cache.IUserCache
}

func NewConversationImplement(
	conversationRepo repositories.IConversationRepository,
	userCache cache.IUserCache,
) *sConversation {
	return &sConversation{
		conversationRepo: conversationRepo,
		userCache:        userCache,
	}
}

func (s *sConversation) GetConversationById(
	ctx context.Context,
	conversationId uuid.UUID,
) (result *common.ConversationResult, err error) {
	conversation, err := s.conversationRepo.GetById(ctx, conversationId)
	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, err
	}

	return mapper.NewConversationResult(conversation), nil
}

func (s *sConversation) CreateConversation(
	ctx context.Context,
	command *conversationCommand.CreateConversationCommand,
) (result *conversationCommand.CreateConversationResult, err error) {
	conversationEntity, err := conversationEntity.NewConversation(command.Name, command.UserIds)
	if err != nil {
		return nil, err
	}

	conversation, err := s.conversationRepo.CreateOne(ctx, conversationEntity)
	if err != nil {
		return nil, err
	}

	return &conversationCommand.CreateConversationResult{
		Conversation: mapper.NewConversationResult(conversation),
	}, nil
}

func (s *sConversation) GetManyConversation(
	ctx context.Context,
	userId uuid.UUID,
	query *conversationQuery.GetManyConversationQuery,
) (result *conversationQuery.GetManyConversationQueryResult, err error) {
	// 1. Get list conversation from db
	conversationEntities, paging, err := s.conversationRepo.GetManyConversation(ctx, userId, query)
	if err != nil {
		return result, err
	}

	// 2. Map
	var wg sync.WaitGroup
	conversationResultsChan := make(chan *common.ConversationWithActiveStatusResult, len(conversationEntities))
	for _, conversation := range conversationEntities {
		wg.Add(1)
		go func(conversation *conversationEntity.Conversation) {
			defer wg.Done()
			userActiveStatus := false
			if conversation.UserID != nil {
				userActiveStatus = s.userCache.IsOnline(ctx, *conversation.UserID)
			}
			conversationResultsChan <- mapper.NewConversationWithActiveStatusResult(conversation, userActiveStatus)
		}(conversation)
	}

	wg.Wait()
	close(conversationResultsChan)

	var conversationResults []*common.ConversationWithActiveStatusResult
	for conversationResult := range conversationResultsChan {
		conversationResults = append(conversationResults, conversationResult)
	}

	return &conversationQuery.GetManyConversationQueryResult{
		Conversation:   conversationResults,
		PagingResponse: paging,
	}, nil
}

func (s *sConversation) DeleteConversationById(
	ctx context.Context,
	command *conversationCommand.DeleteConversationCommand,
) error {
	//1. Find conversation
	conversationFound, err := s.conversationRepo.GetById(ctx, *command.ConversationId)
	if err != nil {
		return err
	}

	if conversationFound == nil {
		return err
	}

	//2. Delete conversation
	if err = s.conversationRepo.DeleteById(ctx, *command.ConversationId); err != nil {
		return err
	}

	return nil
}

func (s *sConversation) UpdateConversationById(
	ctx context.Context,
	command *conversationCommand.UpdateConversationCommand,
) (result *conversationCommand.UpdateConversationCommandResult, err error) {
	conversationFound, err := s.conversationRepo.GetById(ctx, *command.ConversationId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}
	if conversationFound == nil {
		return nil, response.NewDataNotFoundError("conversation not found")
	}

	updateConversationEntity := &conversationEntity.ConversationUpdate{
		Name: command.Name,
	}
	err = updateConversationEntity.ValidateConversationUpdate()
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if command.Image != nil && command.Image.Size > 0 && command.Image.Filename != "" {
		image, err := media.SaveMedia(command.Image)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		_, err = s.conversationRepo.UpdateOne(ctx, *command.ConversationId, &conversationEntity.ConversationUpdate{
			Image: &image,
		})

		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}
	}

	conversationFound, err = s.conversationRepo.UpdateOne(ctx, *command.ConversationId, updateConversationEntity)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}
	return &conversationCommand.UpdateConversationCommandResult{
		Conversation: mapper.NewConversationResult(conversationFound),
	}, nil
}
