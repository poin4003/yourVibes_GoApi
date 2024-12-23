package implement

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	user_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
)

type sUserReport struct {
	userReportRepo user_report_repo.IUserReportRepository
	userRepo       user_report_repo.IUserRepository
	postRepo       user_report_repo.IPostRepository
	commentRepo    user_report_repo.ICommentRepository
}

func NewUserReportImplement(
	userReportRepo user_report_repo.IUserReportRepository,
	userRepo user_report_repo.IUserRepository,
	postRepo user_report_repo.IPostRepository,
	commentRepo user_report_repo.ICommentRepository,
) *sUserReport {
	return &sUserReport{
		userReportRepo: userReportRepo,
		userRepo:       userRepo,
		postRepo:       postRepo,
		commentRepo:    commentRepo,
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
	userReportFound, err := s.userReportRepo.GetById(ctx, command.UserId, command.ReportedUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check if report is already handled
	if userReportFound.Status {
		result.ResultCode = response.ErrCodeReportIsAlreadyHandled
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("You don't need to handle this report again")
	}

	// 3. Update reported user status
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

	// 4. Update reportedUser posts status
	postUpdateEntity := &post_entity.PostUpdate{
		Status: pointer.Ptr(false),
	}
	if err = postUpdateEntity.ValidatePostUpdate(); err != nil {
		return result, err
	}

	conditions := map[string]interface{}{
		"user_id": command.ReportedUserId,
	}

	if err = s.postRepo.UpdateMany(ctx, conditions, postUpdateEntity); err != nil {
		return result, err
	}

	// 5. Update reportedUser comment status
	conditions = map[string]interface{}{
		"user_id": command.ReportedUserId,
	}

	updateData := map[string]interface{}{
		"status": false,
	}

	if err = s.commentRepo.UpdateMany(ctx, conditions, updateData); err != nil {
		return result, err
	}

	// 6. Update report status
	userReportEntity := &user_entity.UserReportUpdate{
		AdminId: pointer.Ptr(command.AdminId),
		Status:  pointer.Ptr(true),
	}

	if err = s.userReportRepo.UpdateMany(ctx, command.ReportedUserId, userReportEntity); err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserReport) DeleteUserReport(
	ctx context.Context,
	command *user_command.DeleteUserReportCommand,
) (result *user_command.DeleteUserReportCommandResult, err error) {
	result = &user_command.DeleteUserReportCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Check exists
	userReportFound, err := s.userReportRepo.GetById(ctx, command.UserId, command.ReportedUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check if report is already handled
	if userReportFound.Status {
		result.ResultCode = response.ErrCodeReportIsAlreadyHandled
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("You can't delete this report, it's already handled")
	}

	// 3. Delete report
	if err = s.userReportRepo.DeleteOne(ctx, command.UserId, command.ReportedUserId); err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserReport) ActivateUserAccount(
	ctx context.Context,
	command *user_command.ActivateUserAccountCommand,
) (result *user_command.ActivateUserAccountCommandResult, err error) {
	result = &user_command.ActivateUserAccountCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Check exists
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check if user is already activate
	if userFound.Status {
		result.ResultCode = response.ErrCodeUserIsAlreadyActivated
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("You don't need to activate this user account")
	}

	// 3. Update reported user status
	reportedUserUpdateEntity := &user_entity.UserUpdate{
		Status: pointer.Ptr(true),
	}

	if err = reportedUserUpdateEntity.ValidateUserUpdate(); err != nil {
		return result, err
	}

	_, err = s.userRepo.UpdateOne(ctx, command.UserId, reportedUserUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 4. Update reportedUser posts status
	postUpdateEntity := &post_entity.PostUpdate{
		Status: pointer.Ptr(true),
	}
	if err = postUpdateEntity.ValidatePostUpdate(); err != nil {
		return result, err
	}

	conditions := map[string]interface{}{
		"user_id": command.UserId,
	}

	if err = s.postRepo.UpdateMany(ctx, conditions, postUpdateEntity); err != nil {
		return result, err
	}

	// 5. Update reportedUser comment status
	conditions = map[string]interface{}{
		"user_id": command.UserId,
	}

	updateData := map[string]interface{}{
		"status": true,
	}

	if err = s.commentRepo.UpdateMany(ctx, conditions, updateData); err != nil {
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
