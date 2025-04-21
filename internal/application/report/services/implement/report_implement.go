package implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"

	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/producer"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/voucher/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/sendto"
	"go.uber.org/zap"

	reportCommand "github.com/poin4003/yourVibes_GoApi/internal/application/report/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/mapper"
	reportQuery "github.com/poin4003/yourVibes_GoApi/internal/application/report/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sReportFactory struct {
	reportRepo            repo.IReportRepository
	voucherRepo           repo.IVoucherRepository
	friendRepo            repo.IFriendRepository
	userCache             cache.IUserCache
	postCache             cache.IPostCache
	commentCache          cache.ICommentCache
	notificationPublisher *producer.NotificationPublisher
}

func NewReportFactoryImplment(
	reportRepo repo.IReportRepository,
	voucherRepo repo.IVoucherRepository,
	friendRepo repo.IFriendRepository,
	userCache cache.IUserCache,
	postCache cache.IPostCache,
	commentCache cache.ICommentCache,
	notificationPublisher *producer.NotificationPublisher,
) *sReportFactory {
	return &sReportFactory{
		reportRepo:            reportRepo,
		voucherRepo:           voucherRepo,
		friendRepo:            friendRepo,
		userCache:             userCache,
		postCache:             postCache,
		commentCache:          commentCache,
		notificationPublisher: notificationPublisher,
	}
}

func (s *sReportFactory) CreateReport(
	ctx context.Context,
	command *reportCommand.CreateReportCommand,
) error {
	var entity interface{}
	var err error

	switch command.Type {
	case consts.USER_REPORT:
		entity, err = reportEntity.NewUserReport(command.Reason, command.Type, command.UserId, command.ReportedId)
		if err != nil {
			return response.NewServerFailedError(err.Error())
		}
	case consts.POST_REPORT:
		entity, err = reportEntity.NewPostReport(command.Reason, command.Type, command.UserId, command.ReportedId)
		if err != nil {
			return response.NewServerFailedError(err.Error())
		}
	case consts.COMMENT_REPORT:
		entity, err = reportEntity.NewCommentReport(command.Reason, command.Type, command.UserId, command.ReportedId)
		if err != nil {
			return response.NewServerFailedError(err.Error())
		}
	default:
		return response.NewValidateError("invalid report type")
	}

	switch e := entity.(type) {
	case *reportEntity.UserReportEntity:
		return s.reportRepo.CreateUserReport(ctx, e)
	case *reportEntity.PostReportEntity:
		return s.reportRepo.CreatePostReport(ctx, e)
	case *reportEntity.CommentReportEntity:
		return s.reportRepo.CreateCommentReport(ctx, e)
	default:
		return response.NewServerFailedError("unsupported report type")
	}
}

func (s *sReportFactory) GetDetailReport(
	ctx context.Context,
	query *reportQuery.GetOneReportQuery,
) (*reportQuery.ReportQueryResult, error) {
	switch query.ReportType {
	case consts.USER_REPORT:
		entity, err := s.reportRepo.GetUserReportById(ctx, query.ReportedId)
		if err != nil {
			return nil, err
		}
		return &reportQuery.ReportQueryResult{
			Type:       consts.USER_REPORT,
			UserReport: mapper.NewUserReportResult(entity),
		}, nil
	case consts.POST_REPORT:
		entity, err := s.reportRepo.GetPostReportById(ctx, query.ReportedId)
		if err != nil {
			return nil, err
		}
		return &reportQuery.ReportQueryResult{
			Type:       consts.POST_REPORT,
			PostReport: mapper.NewPostReportResult(entity),
		}, nil
	case consts.COMMENT_REPORT:
		entity, err := s.reportRepo.GetCommentReportById(ctx, query.ReportedId)
		if err != nil {
			return nil, err
		}
		return &reportQuery.ReportQueryResult{
			Type:          consts.COMMENT_REPORT,
			CommentReport: mapper.NewCommentReportResult(entity),
		}, nil
	default:
		return nil, response.NewValidateError("invalid report type")
	}
}

func (s *sReportFactory) GetManyReport(
	ctx context.Context,
	query *reportQuery.GetManyReportQuery,
) (result *reportQuery.ReportQueryListResult, err error) {
	switch query.ReportType {
	case consts.USER_REPORT:
		entities, paging, err := s.reportRepo.GetManyUserReport(ctx, query)
		if err != nil {
			return nil, err
		}
		var userReportResults []*common.UserReportShortVerResult
		for _, userReportEntity := range entities {
			userReportResult := mapper.NewUserReportShortVerResult(userReportEntity)
			userReportResults = append(userReportResults, userReportResult)
		}
		return &reportQuery.ReportQueryListResult{
			Type:           consts.USER_REPORT,
			UserReports:    userReportResults,
			PagingResponse: paging,
		}, nil
	case consts.POST_REPORT:
		entities, paging, err := s.reportRepo.GetManyPostReport(ctx, query)
		if err != nil {
			return nil, err
		}
		var postReportResults []*common.PostReportShortVerResult
		for _, postReportEntity := range entities {
			postReportResult := mapper.NewPostReportShortVerResult(postReportEntity)
			postReportResults = append(postReportResults, postReportResult)
		}
		return &reportQuery.ReportQueryListResult{
			Type:           consts.POST_REPORT,
			PostReports:    postReportResults,
			PagingResponse: paging,
		}, nil
	case consts.COMMENT_REPORT:
		entities, paging, err := s.reportRepo.GetManyCommentReport(ctx, query)
		if err != nil {
			return nil, err
		}
		var commentReportResults []*common.CommentReportShortVerResult
		for _, commentReportEntity := range entities {
			commentReportResult := mapper.NewCommentReportShortVerResult(commentReportEntity)
			commentReportResults = append(commentReportResults, commentReportResult)
		}
		return &reportQuery.ReportQueryListResult{
			Type:           consts.COMMENT_REPORT,
			CommentReports: commentReportResults,
			PagingResponse: paging,
		}, nil
	default:
		return nil, response.NewValidateError("invalid report type")
	}
}

func (s *sReportFactory) HandleReport(
	ctx context.Context,
	command *reportCommand.HandleReportCommand,
) (err error) {
	switch command.Type {
	case consts.USER_REPORT:
		// Handle user report
		userEntity, err := s.reportRepo.HandleUserReport(ctx, command.ReportId, command.AdminId)
		if err != nil {
			return err
		}

		// Send email for user
		if err = sendto.SendTemplateEmail(
			[]string{userEntity.Email},
			consts.HOST_EMAIL,
			"deactivate_account.html",
			map[string]interface{}{
				"familyname": userEntity.FamilyName,
				"name":       userEntity.Name,
				"email":      userEntity.Email,
			},
			"deactivate account user",
		); err != nil {
			return response.NewServerFailedError(err.Error())
		}

		// Delete user status cache
		go func(userID uuid.UUID) {
			s.userCache.DeleteUserStatus(ctx, userID)
			s.postCache.DeleteRelatedPost(ctx, consts.RK_PERSONAL_POST, userID)
			friendIds, err := s.friendRepo.GetFriendIds(ctx, userID)
			if err != nil {
				global.Logger.Error("Failed to get friendIds", zap.String("user_id", userID.String()), zap.Error(err))
				return
			}
			if len(friendIds) == 0 {
				return
			}
			s.postCache.DeleteFriendFeeds(ctx, consts.RK_USER_FEED, friendIds)
		}(userEntity.ID)

		return nil
	case consts.POST_REPORT:
		// handle post report
		postEntity, err := s.reportRepo.HandlePostReport(ctx, command.ReportId, command.AdminId)
		if err != nil {
			return err
		}

		// Send notification for user
		notification, err := notificationEntity.NewNotification(
			"System",
			consts.AVATAR_URL,
			postEntity.User.ID,
			consts.DEACTIVATE_POST,
			postEntity.ID.String(),
			"your post has been deactivate",
		)
		if err != nil {
			global.Logger.Error("Failed to create notification entity", zap.Error(err))
			return nil
		}

		// Publish to RabbitMQ to handle Notification
		notiMsg := mapper.NewNotificationResult(notification)
		if err = s.notificationPublisher.PublishNotification(ctx, notiMsg, "notification.single.db_websocket"); err != nil {
			global.Logger.Error("Failed to publish notification for friend", zap.Error(err))
		}

		// Delete cache
		go func(postId, userId uuid.UUID) {
			s.postCache.DeletePost(ctx, postId)
			s.postCache.DeleteRelatedPost(ctx, consts.RK_PERSONAL_POST, userId)
			friendIds, err := s.friendRepo.GetFriendIds(ctx, userId)
			if err != nil {
				global.Logger.Error("Failed to get friendIds", zap.String("user_id", userId.String()), zap.Error(err))
				return
			}
			if len(friendIds) == 0 {
				return
			}
			s.postCache.DeleteFriendFeeds(ctx, consts.RK_USER_FEED, friendIds)
		}(postEntity.ID, postEntity.UserId)

		return nil
	case consts.COMMENT_REPORT:
		// Handle comment report
		commentEntity, err := s.reportRepo.HandleCommentReport(ctx, command.ReportId, command.AdminId)
		if err != nil {
			return err
		}

		// Send notification for user
		notification, err := notificationEntity.NewNotification(
			"System",
			consts.AVATAR_URL,
			commentEntity.User.ID,
			consts.DEACTIVATE_COMMENT,
			commentEntity.ID.String(),
			"your comment has been deactivated",
		)
		if err != nil {
			global.Logger.Error("Failed to create notification entity", zap.Error(err))
			return nil
		}

		// Publish to RabbitMQ to handle Notification
		notiMsg := mapper.NewNotificationResult(notification)
		if err = s.notificationPublisher.PublishNotification(ctx, notiMsg, "notification.single.db_websocket"); err != nil {
			global.Logger.Error("Failed to publish notification for friend", zap.Error(err))
		}

		//s.commentCache.DeleteComment(ctx, commentEntity.ID)

		return nil
	default:
		return response.NewValidateError("invalid report type")
	}
}

func (s *sReportFactory) DeleteReport(
	ctx context.Context,
	command *reportCommand.DeleteReportCommand,
) (err error) {
	if err = s.reportRepo.DeleteReportById(ctx, command.ReportId); err != nil {
		return err
	}
	return nil
}

func (s *sReportFactory) Activate(
	ctx context.Context,
	command *reportCommand.ActivateCommand,
) (err error) {
	switch command.Type {
	case consts.USER_REPORT:
		// Activate user
		userEntity, err := s.reportRepo.ActivateUser(ctx, command.ReportId)
		if err != nil {
			return err
		}

		// Send email for user
		if err = sendto.SendTemplateEmail(
			[]string{userEntity.Email},
			consts.HOST_EMAIL,
			"activate_account.html",
			map[string]interface{}{
				"familyname": userEntity.FamilyName,
				"name":       userEntity.Name,
				"email":      userEntity.Email,
			},
			"activate account for user",
		); err != nil {
			return response.NewServerFailedError(err.Error())
		}

		// Delete user status cache
		go func(userID uuid.UUID) {
			s.userCache.DeleteUserStatus(ctx, userID)
			s.postCache.DeleteRelatedPost(ctx, consts.RK_PERSONAL_POST, userID)
			friendIds, err := s.friendRepo.GetFriendIds(ctx, userID)
			if err != nil {
				global.Logger.Error("Failed to get friendIds", zap.String("user_id", userID.String()), zap.Error(err))
				return
			}
			if len(friendIds) == 0 {
				return
			}
			s.postCache.DeleteFriendFeeds(ctx, consts.RK_USER_FEED, friendIds)
		}(userEntity.ID)

		return nil
	case consts.POST_REPORT:
		// Activate post
		postEntity, err := s.reportRepo.ActivatePost(ctx, command.ReportId)
		if err != nil {
			return err
		}

		// Check if post is advertise or was advertises
		content := "your post has been activated"
		if postEntity.IsAdvertisement != consts.NOT_ADVERTISE {
			voucherEntity, err := entities.NewVoucherBySystem(
				"yourvibes",
				"yourvibes voucher after admin deactivate mistake",
				1,
				20,
				consts.PERCENTAGE,
			)
			if err != nil {
				return response.NewServerFailedError(err.Error())
			}

			if err = s.voucherRepo.CreateVoucher(ctx, voucherEntity); err != nil {
				return err
			}

			content = "your post has been activated, your voucher is " + voucherEntity.Code
		}

		// Send notification for user
		notification, err := notificationEntity.NewNotification(
			"System",
			consts.AVATAR_URL,
			postEntity.User.ID,
			consts.ACTIVATE_POST,
			postEntity.ID.String(),
			content,
		)
		if err != nil {
			global.Logger.Error("Failed to create notification entity", zap.Error(err))
			return nil
		}

		// Publish to RabbitMQ to handle Notification
		notiMsg := mapper.NewNotificationResult(notification)
		if err = s.notificationPublisher.PublishNotification(ctx, notiMsg, "notification.single.db_websocket"); err != nil {
			global.Logger.Error("Failed to publish notification for friend", zap.Error(err))
		}

		go func(postId, userId uuid.UUID) {
			s.postCache.DeletePost(ctx, postId)
			s.postCache.DeleteRelatedPost(ctx, consts.RK_PERSONAL_POST, userId)
			friendIds, err := s.friendRepo.GetFriendIds(ctx, userId)
			if err != nil {
				global.Logger.Error("Failed to get friendIds", zap.String("user_id", userId.String()), zap.Error(err))
				return
			}
			if len(friendIds) == 0 {
				return
			}
			s.postCache.DeleteFriendFeeds(ctx, consts.RK_USER_FEED, friendIds)
		}(postEntity.ID, postEntity.UserId)

		return nil
	case consts.COMMENT_REPORT:
		commentEntity, err := s.reportRepo.ActivateComment(ctx, command.ReportId)
		if err != nil {
			return err
		}

		// Send notification for user
		notification, err := notificationEntity.NewNotification(
			commentEntity.User.FamilyName+" "+commentEntity.User.Name,
			commentEntity.User.AvatarUrl,
			commentEntity.User.ID,
			consts.ACTIVATE_COMMENT,
			commentEntity.ID.String(),
			"your comment has been activated",
		)
		if err != nil {
			global.Logger.Error("Failed to create notification entity", zap.Error(err))
			return nil
		}

		// Publish to RabbitMQ to handle Notification
		notiMsg := mapper.NewNotificationResult(notification)
		if err = s.notificationPublisher.PublishNotification(ctx, notiMsg, "notification.single.db_websocket"); err != nil {
			global.Logger.Error("Failed to publish notification for friend", zap.Error(err))
		}

		return nil
	default:
		return response.NewValidateError("invalid report type")
	}
}
