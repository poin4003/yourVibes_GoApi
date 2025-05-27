package implement

import (
	"context"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/producer"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	conversationDetailCommand "github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/mapper"
	conversationDetailQuery "github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"
)

type sConversationDetail struct {
	conversationRepo       repositories.IConversationRepository
	messageRepo            repositories.IMessageRepository
	conversationDetailRepo repositories.IConversationDetailRepository
	userRepo               repositories.IUserRepository
	messagePublisher       *producer.MessagePublisher
}

func NewConversationDetailImplement(
	conversationRepo repositories.IConversationRepository,
	messageRepo repositories.IMessageRepository,
	conversationDetailRepo repositories.IConversationDetailRepository,
	userRepo repositories.IUserRepository,
	messagePublisher *producer.MessagePublisher,
) *sConversationDetail {
	return &sConversationDetail{
		conversationRepo:       conversationRepo,
		messageRepo:            messageRepo,
		conversationDetailRepo: conversationDetailRepo,
		userRepo:               userRepo,
		messagePublisher:       messagePublisher,
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
	// 1. Check conversation exists
	conversation, err := s.conversationRepo.GetById(ctx, command.ConversationId)
	if err != nil {
		return err
	}

	if conversation == nil {
		return err
	}

	// 2. Create new member
	newConversationDetail, _ := entities.NewConversationDetail(
		command.UserId,
		command.ConversationId,
		consts.CONVERSATION_MEMBER,
	)

	member, err := s.conversationDetailRepo.CreateOne(ctx, newConversationDetail)
	if err != nil {
		return err
	}

	UserCommand := conversationDetailCommand.UserCommand{
		ID:         member.UserId.String(),
		AvatarUrl:  member.User.AvatarUrl,
		Name:       member.User.Name,
		FamilyName: member.User.FamilyName,
	}

	createMessageCommand := &conversationDetailCommand.CreateMessageCommand{
		ConversationId: command.ConversationId,
		UserId:         command.UserId,
		User:           UserCommand,
		ParentId:       nil,
		ParentContent:  nil,
		Content:        fmt.Sprintf("%s %s join the conversation", member.User.FamilyName, member.User.Name),
		CreatedAt:      time.Now(),
	}

	// 3. Publish to rabbitmq (push websocket)
	go func(createMessageCommand *conversationDetailCommand.CreateMessageCommand) {
		if err = s.messagePublisher.PublishMessage(ctx, createMessageCommand); err != nil {
			global.Logger.Error("Failed to publish message to rabbitmq", zap.Error(err))
		}
	}(createMessageCommand)

	newMessage, _ := entities.NewMessage(
		command.UserId,
		command.ConversationId,
		nil,
		&createMessageCommand.Content,
	)

	// 4. Save message
	if err = s.messageRepo.CreateOne(ctx, newMessage); err != nil {
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
	// 1. Delete from db
	if err := s.conversationDetailRepo.DeleteById(ctx, *command.UserId, command.AuthenticatedUserId, *command.ConversationId); err != nil {
		return err
	}

	// 2. Get user info
	userFound, err := s.userRepo.GetById(ctx, *command.UserId)
	if err != nil {
		return err
	}

	UserCommand := conversationDetailCommand.UserCommand{
		ID:         userFound.ID.String(),
		AvatarUrl:  userFound.AvatarUrl,
		Name:       userFound.Name,
		FamilyName: userFound.FamilyName,
	}

	createMessageCommand := &conversationDetailCommand.CreateMessageCommand{
		ConversationId: *command.ConversationId,
		UserId:         *command.UserId,
		ParentId:       nil,
		ParentContent:  nil,
		User:           UserCommand,
		Content:        fmt.Sprintf("%s %s left the conversation", userFound.FamilyName, userFound.Name),
		CreatedAt:      time.Now(),
	}
	// 3. Publish to rabbitmq (push websocket)
	go func(createMessageCommand *conversationDetailCommand.CreateMessageCommand) {
		if err = s.messagePublisher.PublishMessage(ctx, createMessageCommand); err != nil {
			global.Logger.Error("Failed to push message to rabbitmq", zap.Error(err))
		}
	}(createMessageCommand)

	newMessage, _ := entities.NewMessage(
		*command.UserId,
		*command.ConversationId,
		nil,
		&createMessageCommand.Content,
	)

	// 4. Save message
	if err = s.messageRepo.CreateOne(ctx, newMessage); err != nil {
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
	// 1. Check conversation exist
	conversation, err := s.conversationRepo.GetById(ctx, command.ConversationId)
	if err != nil {
		return err
	}
	if conversation == nil {
		return err
	}
	// 2. Map to entity to create
	var conversationDetails []*entities.ConversationDetail
	for _, userId := range command.UserIds {
		newConversationDetail, _ := entities.NewConversationDetail(userId, command.ConversationId, consts.CONVERSATION_MEMBER)
		conversationDetails = append(conversationDetails, newConversationDetail)
	}

	// 3. Create many members for conversation
	members, err := s.conversationDetailRepo.CreateMany(ctx, conversationDetails)
	if err != nil {
		return err
	}

	// 4. Publish to rabbitmq (push websocket)
	var messages []*entities.Message
	for _, member := range members {
		content := fmt.Sprintf("%s %s joined the conversation", member.User.FamilyName, member.User.Name)

		createMessageCommand := &conversationDetailCommand.CreateMessageCommand{
			ConversationId: command.ConversationId,
			UserId:         member.UserId,
			ParentId:       nil,
			ParentContent:  nil,
			Content:        content,
			CreatedAt:      time.Now(),
		}
		go func(cmd *conversationDetailCommand.CreateMessageCommand) {
			if err := s.messagePublisher.PublishMessage(ctx, cmd); err != nil {
				global.Logger.Error("Failed to publish message to RabbitMQ", zap.Error(err))
			}
		}(createMessageCommand)

		newMessage, _ := entities.NewMessage(
			member.UserId,
			member.ConversationId,
			nil,
			&content,
		)
		messages = append(messages, newMessage)
	}

	// 5. Save many message
	if err = s.messageRepo.CreateMany(ctx, messages); err != nil {
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
