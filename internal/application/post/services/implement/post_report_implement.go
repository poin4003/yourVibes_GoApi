package implement

import (
	"context"
	"fmt"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	post_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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
