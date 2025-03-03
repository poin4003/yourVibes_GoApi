package implement

import (
	"context"

	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/truncate"

	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	postReportRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
)

type sPostReport struct {
	postReportRepo   postReportRepo.IPostReportRepository
	postRepo         postReportRepo.IPostRepository
	notificationRepo postReportRepo.INotificationRepository
}

func NewPostReportImplement(
	postReportRepo postReportRepo.IPostReportRepository,
	postRepo postReportRepo.IPostRepository,
	notificationRepo postReportRepo.INotificationRepository,
) *sPostReport {
	return &sPostReport{
		postReportRepo:   postReportRepo,
		postRepo:         postRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sPostReport) CreatePostReport(
	ctx context.Context,
	command *postCommand.CreateReportPostCommand,
) (result *postCommand.CreateReportPostCommandResult, err error) {
	// 1. Check report exist
	postReportCheck, err := s.postReportRepo.CheckExist(ctx, command.UserId, command.ReportedPostId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 2. Return if report has already exists
	if postReportCheck {
		return nil, response.NewCustomError(response.ErrCodePostReportHasAlreadyExist)
	}

	// 3. Create report
	postReportEntity, err := postEntity.NewPostReport(
		command.UserId,
		command.ReportedPostId,
		command.Reason,
	)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	userReport, err := s.postReportRepo.CreateOne(ctx, postReportEntity)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 4. Map to result
	return &postCommand.CreateReportPostCommandResult{
		PostReport: mapper.NewPostReportResult(userReport),
	}, nil
}

func (s *sPostReport) HandlePostReport(
	ctx context.Context,
	command *postCommand.HandlePostReportCommand,
) (err error) {
	// 1. Check exist
	postReportFound, err := s.postReportRepo.GetById(ctx, command.UserId, command.ReportedPostId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if postReportFound == nil {
		return response.NewDataNotFoundError("post report not found")
	}

	// 2. Check if report is already handled
	if postReportFound.Status {
		return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
	}

	// 3. Update reported post status
	reportedPostUpdateEntity := &postEntity.PostUpdate{
		Status: pointer.Ptr(false),
	}

	if err = reportedPostUpdateEntity.ValidatePostUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	postUpdated, err := s.postRepo.UpdateOne(ctx, command.ReportedPostId, reportedPostUpdateEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Update report status
	postReportEntity := &postEntity.PostReportUpdate{
		AdminId: pointer.Ptr(command.AdminId),
		Status:  pointer.Ptr(true),
	}

	if err = s.postReportRepo.UpdateMany(ctx, command.ReportedPostId, postReportEntity); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 5. Create notification for user
	content := "Your post has been blocked for violating our policies: " + truncate.TruncateContent(postUpdated.Content, 20)
	notification, err := notificationEntity.NewNotification(
		postUpdated.User.FamilyName+" "+postUpdated.User.Name,
		postUpdated.User.AvatarUrl,
		postUpdated.User.ID,
		consts.DEACTIVATE_POST,
		postUpdated.ID.String(),
		content,
	)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.notificationRepo.CreateOne(ctx, notification)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sPostReport) DeletePostReport(
	ctx context.Context,
	command *postCommand.DeletePostReportCommand,
) (err error) {
	// 1. Check exist
	postReportFound, err := s.postReportRepo.GetById(ctx, command.UserId, command.ReportedPostId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if postReportFound == nil {
		return response.NewDataNotFoundError("post report not found")
	}

	// 2. Check if report is already handled
	if postReportFound.Status {
		return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
	}

	// 3. Delete report
	if err = s.postReportRepo.DeleteOne(ctx, command.UserId, command.ReportedPostId); err != nil {
		return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
	}

	return nil
}

func (s *sPostReport) ActivatePost(
	ctx context.Context,
	command *postCommand.ActivatePostCommand,
) (err error) {
	// 1. Check exist
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if postFound == nil {
		return response.NewDataNotFoundError("post not found")
	}

	// 2. Check if post already activate
	if postFound.Status {
		return response.NewCustomError(response.ErrCodePostIsAlreadyActivated)
	}

	// 3. Update reported post status
	reportedPostUpdateEntity := &postEntity.PostUpdate{
		Status: pointer.Ptr(true),
	}

	if err = reportedPostUpdateEntity.ValidatePostUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	postUpdated, err := s.postRepo.UpdateOne(ctx, command.PostId, reportedPostUpdateEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Delete report
	if err = s.postReportRepo.DeleteByPostId(ctx, command.PostId); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 5. Create notification for user
	content := "Your post is activated: " + truncate.TruncateContent(postUpdated.Content, 20)
	notification, err := notificationEntity.NewNotification(
		postUpdated.User.FamilyName+" "+postUpdated.User.Name,
		postUpdated.User.AvatarUrl,
		postUpdated.User.ID,
		consts.ACTIVATE_POST,
		postUpdated.ID.String(),
		content,
	)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.notificationRepo.CreateOne(ctx, notification)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sPostReport) GetDetailPostReport(
	ctx context.Context,
	query *postQuery.GetOnePostReportQuery,
) (result *postQuery.PostReportQueryResult, err error) {
	// 1. Get post report detail
	postReportEntity, err := s.postReportRepo.GetById(ctx, query.UserId, query.ReportedPostId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if postReportEntity == nil {
		return nil, response.NewDataNotFoundError("post report not found")
	}

	// 2. Map to result

	return &postQuery.PostReportQueryResult{
		PostReport: mapper.NewPostReportResult(postReportEntity),
	}, nil
}

func (s *sPostReport) GetManyPostReport(
	ctx context.Context,
	query *postQuery.GetManyPostReportQuery,
) (result *postQuery.PostReportQueryListResult, err error) {
	// 1. Get list of post report
	postReportEntities, paging, err := s.postReportRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	// 2. Map to result
	var postReportResults []*common.PostReportShortVerResult
	for _, postReportEntity := range postReportEntities {
		postReportResult := mapper.NewPostReportShortVerResult(postReportEntity)
		postReportResults = append(postReportResults, postReportResult)
	}

	return &postQuery.PostReportQueryListResult{
		PostReports:    postReportResults,
		PagingResponse: paging,
	}, nil
}
