package implement

import (
	"context"
	"fmt"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	user_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type sUserReport struct {
	userReportRepo user_report_repo.IUserReportRepository
}

func NewUserReportImplement(
	userReportRepo user_report_repo.IUserReportRepository,
) *sUserReport {
	return &sUserReport{
		userReportRepo: userReportRepo,
	}
}

func (s *sUserReport) CreateUserReport(
	ctx context.Context,
	command *user_command.CreateReportUserCommand,
) (result *user_command.CreateReportUserCommandResult, err error) {
	result = &user_command.CreateReportUserCommandResult{}
	result.UserReport = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Check report exist
	userReportCheck, err := s.userReportRepo.CheckExist(ctx, command.UserId, command.ReportedUserId)
	if err != nil {
		return result, err
	}

	// 2. Return if report has already exist
	if userReportCheck {
		result.ResultCode = response.ErrCodeUserReportHasAlreadyExist
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("user report already exist")
	}

	// 3. Create report
	userReportEntity, err := user_entity.NewUserReport(
		command.UserId,
		command.ReportedUserId,
		command.Reason,
	)
	if err != nil {
		return result, err
	}

	userReport, err := s.userReportRepo.CreateOne(ctx, userReportEntity)
	if err != nil {
		return result, err
	}

	// 4. Map to result
	result.UserReport = mapper.NewUserReportResult(userReport)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
