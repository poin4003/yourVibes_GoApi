package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	conversationDetailCommand "github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/mapper"
	conversationDetailQuery "github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	messageRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"
)

type sConversationDetail struct {
	conversationRepo       messageRepo.IConversationRepository
	messageRepo            messageRepo.IMessageRepository
	conversationDetailRepo messageRepo.IConversationDetailRepository
}

func NewConversationDetailImplement(
	conversationRepo messageRepo.IConversationRepository,
	messageRepo messageRepo.IMessageRepository,
	conversationDetailRepo messageRepo.IConversationDetailRepository,
) *sConversationDetail {
	return &sConversationDetail{
		conversationRepo:       conversationRepo,
		messageRepo:            messageRepo,
		conversationDetailRepo: conversationDetailRepo,
	}
}

func (s *sConversationDetail) GetConversationDetailById(
	ctx context.Context,
	userId uuid.UUID,
	conversationId uuid.UUID,
) (result *common.ConversationDetailResult, err error) {
	conversationDetail, err := s.conversationDetailRepo.GetById(ctx, userId, conversationId)
	if err != nil {
		return nil, err
	}

	if conversationDetail == nil {
		return nil, err
	}

	return mapper.NewConversationDetailResult(conversationDetail), nil
}

func (s *sConversationDetail) CreateConversationDetail(
	ctx context.Context,
	command *conversationDetailCommand.CreateConversationDetailCommand,
) error {
	conversation, err := s.conversationRepo.GetById(ctx, command.ConversationId)
	if err != nil {
		return err
	}

	if conversation == nil {
		return err
	}

	newConversationDetail, _ := entities.NewConversationDetail(
		command.UserId,
		command.ConversationId,
		consts.CONVERSATION_MEMBER,
	)

	_, err = s.conversationDetailRepo.CreateOne(ctx, newConversationDetail)
	if err != nil {
		return err
	}

	return nil
}

func (s *sConversationDetail) GetConversationDetailByConversationId(
	ctx context.Context,
	query *conversationDetailQuery.GetConversationDetailQuery,
) (result *conversationDetailQuery.GetConversationDetailResult, err error) {
	conversationDetailEntities, paging, err := s.conversationDetailRepo.GetConversationDetailByConversationId(ctx, query)
	if err != nil {
		return result, err
	}

	var conversationDetailResults []*common.ConversationDetailResult
	for _, conversationDetailEntity := range conversationDetailEntities {
		conversationDetailResults = append(conversationDetailResults, mapper.NewConversationDetailResult(conversationDetailEntity))
	}

	return &conversationDetailQuery.GetConversationDetailResult{
		ConversationDetail: conversationDetailResults,
		PagingResponse:     paging,
	}, nil
}

func (s *sConversationDetail) DeleteConversationDetailById(
	ctx context.Context,
	command *conversationDetailCommand.DeleteConversationDetailCommand,
) error {
	if err := s.conversationDetailRepo.DeleteById(ctx, *command.UserId, command.AuthenticatedUserId, *command.ConversationId); err != nil {
		return err
	}

	return nil
}

func (s *sConversationDetail) UpdateOneStatusConversationDetail(
	ctx context.Context,
	command *command.UpdateOneStatusConversationDetailCommand,
) (err error) {
	notificationUpdateEntity := &entities.ConversationDetailUpdate{
		LastMessStatus: pointer.Ptr(false),
	}

	newConversationUpdateEntity, _ := entities.NewConversationDetailUpdate(notificationUpdateEntity)

	_, err = s.conversationDetailRepo.UpdateOneStatus(ctx, command.UserId, command.ConversationId, newConversationUpdateEntity)

	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sConversationDetail) CreateManyConversationDetail(
	ctx context.Context,
	command *conversationDetailCommand.CreateManyConversationDetailCommand,
) error {
	conversation, err := s.conversationRepo.GetById(ctx, command.ConversationId)
	if err != nil {
		return err
	}
	if conversation == nil {
		return err
	}
	var conversationDetails []*entities.ConversationDetail
	for _, userId := range command.UserIds {
		newConversationDetail, _ := entities.NewConversationDetail(userId, command.ConversationId, consts.CONVERSATION_MEMBER)
		conversationDetails = append(conversationDetails, newConversationDetail)
	}
	err = s.conversationDetailRepo.CreateMany(ctx, conversationDetails)
	if err != nil {
		return err
	}

	return nil
}

func (s *sConversationDetail) TransferOwnerRole(
	ctx context.Context,
	command *command.TransferOwnerRoleCommand,
) (err error) {
	if err = s.conversationDetailRepo.TransferOwnerRole(ctx,
		command.UserId, command.AuthenticatedUserId, command.ConversationId,
	); err != nil {
		return err
	}

	return nil
}
