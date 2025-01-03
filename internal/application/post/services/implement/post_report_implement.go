package implement

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	postReportRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
)

type sPostReport struct {
	postReportRepo postReportRepo.IPostReportRepository
	postRepo       postReportRepo.IPostRepository
}

func NewPostReportImplement(
	postReportRepo postReportRepo.IPostReportRepository,
	postRepo postReportRepo.IPostRepository,
) *sPostReport {
	return &sPostReport{
		postReportRepo: postReportRepo,
		postRepo:       postRepo,
	}
}

func (s *sPostReport) CreatePostReport(
	ctx context.Context,
	command *postCommand.CreateReportPostCommand,
) (result *postCommand.CreateReportPostCommandResult, err error) {
	result = &postCommand.CreateReportPostCommandResult{}
	result.PostReport = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Check report exist
	postReportCheck, err := s.postReportRepo.CheckExist(ctx, command.UserId, command.ReportedPostId)
	if err != nil {
		return result, err
	}

	// 2. Return if report has already exists
	if postReportCheck {
		result.ResultCode = response.ErrCodePostReportHasAlreadyExist
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("post report already exist")
	}

	// 3. Create report
	postReportEntity, err := postEntity.NewPostReport(
		command.UserId,
		command.ReportedPostId,
		command.Reason,
	)
	if err != nil {
		return result, err
	}

	userReport, err := s.postReportRepo.CreateOne(ctx, postReportEntity)
	if err != nil {
		return result, err
	}

	// 4. Map to result
	result.PostReport = mapper.NewPostReportResult(userReport)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostReport) HandlePostReport(
	ctx context.Context,
	command *postCommand.HandlePostReportCommand,
) (result *postCommand.HandlePostReportCommandResult, err error) {
	result = &postCommand.HandlePostReportCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Check exist
	postReportFound, err := s.postReportRepo.GetById(ctx, command.UserId, command.ReportedPostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check if report is already handled
	if postReportFound.Status {
		result.ResultCode = response.ErrCodeReportIsAlreadyHandled
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you don't need to handle this report again")
	}

	// 3. Update reported post status
	reportedPostUpdateEntity := &postEntity.PostUpdate{
		Status: pointer.Ptr(false),
	}

	if err = reportedPostUpdateEntity.ValidatePostUpdate(); err != nil {
		return result, err
	}

	_, err = s.postRepo.UpdateOne(ctx, command.ReportedPostId, reportedPostUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("post report not found")
		}
		return result, err
	}

	// 4. Update report status
	postReportEntity := &postEntity.PostReportUpdate{
		AdminId: pointer.Ptr(command.AdminId),
		Status:  pointer.Ptr(true),
	}

	if err = s.postReportRepo.UpdateMany(ctx, command.ReportedPostId, postReportEntity); err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostReport) DeletePostReport(
	ctx context.Context,
	command *postCommand.DeletePostReportCommand,
) (result *postCommand.DeletePostReportCommandResult, err error) {
	result = &postCommand.DeletePostReportCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Check exist
	postReportFound, err := s.postReportRepo.GetById(ctx, command.UserId, command.ReportedPostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check if report is already handled
	if postReportFound.Status {
		result.ResultCode = response.ErrCodeReportIsAlreadyHandled
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you can't delete this report, it's already handled")
	}

	// 3. Delete report
	if err = s.postReportRepo.DeleteOne(ctx, command.UserId, command.ReportedPostId); err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostReport) ActivatePost(
	ctx context.Context,
	command *postCommand.ActivatePostCommand,
) (result *postCommand.ActivatePostCommandResult, err error) {
	result = &postCommand.ActivatePostCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Check exist
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check if post already activate
	if postFound.Status {
		result.ResultCode = response.ErrCodePostIsAlreadyActivated
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you don't need to activate this post")
	}

	// 3. Update reported post status
	reportedPostUpdateEntity := &postEntity.PostUpdate{
		Status: pointer.Ptr(true),
	}

	if err = reportedPostUpdateEntity.ValidatePostUpdate(); err != nil {
		return result, err
	}

	_, err = s.postRepo.UpdateOne(ctx, command.PostId, reportedPostUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("post report not found")
		}
		return result, err
	}

	// 4. Delete report
	if err = s.postReportRepo.DeleteByPostId(ctx, command.PostId); err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostReport) GetDetailPostReport(
	ctx context.Context,
	query *postQuery.GetOnePostReportQuery,
) (result *postQuery.PostReportQueryResult, err error) {
	result = &postQuery.PostReportQueryResult{}
	result.PostReport = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Get post report detail
	postReportEntity, err := s.postReportRepo.GetById(ctx, query.UserId, query.ReportedPostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Map to result
	result.PostReport = mapper.NewPostReportResult(postReportEntity)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostReport) GetManyPostReport(
	ctx context.Context,
	query *postQuery.GetManyPostReportQuery,
) (result *postQuery.PostReportQueryListResult, err error) {
	result = &postQuery.PostReportQueryListResult{}
	result.PostReports = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	result.PagingResponse = nil
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

	result.PostReports = postReportResults
	result.PagingResponse = paging
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
