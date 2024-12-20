package implement

import (
	"context"
	"errors"
	"fmt"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	user_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
	"net/http"
)

type sUserReport struct {
	userReportRepo user_report_repo.IUserReportRepository
	userRepo       user_report_repo.IUserRepository
}

func NewUserReportImplement(
	userReportRepo user_report_repo.IUserReportRepository,
	userRepo user_report_repo.IUserRepository,
) *sUserReport {
	return &sUserReport{
		userReportRepo: userReportRepo,
		userRepo:       userRepo,
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

func (s *sUserReport) HandleUserReport(
	ctx context.Context,
	command *user_command.HandleUserReportCommand,
) (result *user_command.HandleUserReportCommandResult, err error) {
	result = &user_command.HandleUserReportCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Check exists
	userReportFound, err := s.userReportRepo.CheckExist(ctx, command.UserId, command.ReportedUserId)
	if err != nil {
		return result, err
	}

	if !userReportFound {
		result.ResultCode = response.ErrDataNotFound
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("user report not found")
	}

	// 2. Update reported user status
	reportedUserUpdateEntity := &user_entity.UserUpdate{
		Status: pointer.Ptr(false),
	}

	if err = reportedUserUpdateEntity.ValidateUserUpdate(); err != nil {
		return result, err
	}

	_, err = s.userRepo.UpdateOne(ctx, command.ReportedUserId, reportedUserUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 3. Update report status
	userReportEntity := &user_entity.UserReportUpdate{
		AdminId: pointer.Ptr(command.AdminId),
		Status:  pointer.Ptr(false),
	}

	if err = s.userReportRepo.UpdateMany(ctx, command.ReportedUserId, userReportEntity); err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserReport) GetDetailUserReport(
	ctx context.Context,
	query *user_query.GetOneUserReportQuery,
) (result *user_query.UserReportQueryResult, err error) {
	result = &user_query.UserReportQueryResult{}
	result.UserReport = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Get user report detail
	userReportEntity, err := s.userReportRepo.GetById(ctx, query.UserId, query.ReportedUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Map to result
	result.UserReport = mapper.NewUserReportResult(userReportEntity)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserReport) GetManyUserReport(
	ctx context.Context,
	query *user_query.GetManyUserReportQuery,
) (result *user_query.UserReportQueryListResult, err error) {
	result = &user_query.UserReportQueryListResult{}
	result.UserReports = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	result.PagingResponse = nil
	// 1. Get list of user report
	userReportEntities, paging, err := s.userReportRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	// 2. Map to result
	var userReportResults []*common.UserReportShortVerResult
	for _, userReportEntity := range userReportEntities {
		userReportResult := mapper.NewUserReportShortVerResult(userReportEntity)
		userReportResults = append(userReportResults, userReportResult)
	}

	result.UserReports = userReportResults
	result.PagingResponse = paging
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
