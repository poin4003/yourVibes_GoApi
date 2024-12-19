package implement

import (
	"context"
	"errors"
	"fmt"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	post_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
)

type sPostReport struct {
	postReportRepo post_report_repo.IPostReportRepository
}

func NewPostReportImplement(
	postReportRepo post_report_repo.IPostReportRepository,
) *sPostReport {
	return &sPostReport{
		postReportRepo: postReportRepo,
	}
}

func (s *sPostReport) CreatePostReport(
	ctx context.Context,
	command *post_command.CreateReportPostCommand,
) (result *post_command.CreateReportPostCommandResult, err error) {
	result = &post_command.CreateReportPostCommandResult{}
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
	postReportEntity, err := post_entity.NewPostReport(
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
	command *post_command.HandlePostReportCommand,
) (result *post_command.HandlePostReportCommandResult, err error) {
	return nil, nil
}

func (s *sPostReport) GetDetailPostReport(
	ctx context.Context,
	query *post_query.GetOnePostReportQuery,
) (result *post_query.PostReportQueryResult, err error) {
	result = &post_query.PostReportQueryResult{}
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
	query *post_query.GetManyPostReportQuery,
) (result *post_query.PostReportQueryListResult, err error) {
	result = &post_query.PostReportQueryListResult{}
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
