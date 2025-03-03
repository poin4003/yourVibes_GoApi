package implement

import (
	"context"

	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/sendto"

	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	userQuery "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	userReportRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
)

type sUserReport struct {
	userReportRepo userReportRepo.IUserReportRepository
	userRepo       userReportRepo.IUserRepository
	postRepo       userReportRepo.IPostRepository
	commentRepo    userReportRepo.ICommentRepository
}

func NewUserReportImplement(
	userReportRepo userReportRepo.IUserReportRepository,
	userRepo userReportRepo.IUserRepository,
	postRepo userReportRepo.IPostRepository,
	commentRepo userReportRepo.ICommentRepository,
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
	command *userCommand.CreateReportUserCommand,
) (result *userCommand.CreateReportUserCommandResult, err error) {
	// 1. Check report exist
	userReportCheck, err := s.userReportRepo.CheckExist(ctx, command.UserId, command.ReportedUserId)
	if err != nil {
		return nil, response.NewCustomError(response.ErrDataNotFound)
	}

	// 2. Return if report has already exist
	if userReportCheck {
		return nil, response.NewCustomError(response.ErrCodeUserReportHasAlreadyExist)
	}

	// 3. Create report
	userReportEntity, err := userEntity.NewUserReport(
		command.UserId,
		command.ReportedUserId,
		command.Reason,
	)
	if err != nil {
		return nil, response.NewCustomError(response.ErrServerFailed)
	}

	userReport, err := s.userReportRepo.CreateOne(ctx, userReportEntity)
	if err != nil {
		return nil, response.NewCustomError(response.ErrServerFailed)
	}

	// 4. Map to result
	return &userCommand.CreateReportUserCommandResult{
		UserReport: mapper.NewUserReportResult(userReport),
	}, nil
}

func (s *sUserReport) HandleUserReport(
	ctx context.Context,
	command *userCommand.HandleUserReportCommand,
) (err error) {
	// 1. Check exists
	userReportFound, err := s.userReportRepo.GetById(ctx, command.UserId, command.ReportedUserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if userReportFound == nil {
		return response.NewDataNotFoundError("user report not found")
	}

	// 2. Check if report is already handled
	if userReportFound.Status {
		return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
	}

	// 3. Update reported user status
	reportedUserUpdateEntity := &userEntity.UserUpdate{
		Status: pointer.Ptr(false),
	}

	if err = reportedUserUpdateEntity.ValidateUserUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	userUpdated, err := s.userRepo.UpdateOne(ctx, command.ReportedUserId, reportedUserUpdateEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Update reportedUser posts status
	postUpdateEntity := &postEntity.PostUpdate{
		Status: pointer.Ptr(false),
	}
	if err = postUpdateEntity.ValidatePostUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	conditions := map[string]interface{}{
		"user_id": command.ReportedUserId,
	}

	if err = s.postRepo.UpdateMany(ctx, conditions, postUpdateEntity); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 5. Update reportedUser comment status
	conditions = map[string]interface{}{
		"user_id": command.ReportedUserId,
	}

	updateData := map[string]interface{}{
		"status": false,
	}

	if err = s.commentRepo.UpdateMany(ctx, conditions, updateData); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 6. Update report status
	userReportEntity := &userEntity.UserReportUpdate{
		AdminId: pointer.Ptr(command.AdminId),
		Status:  pointer.Ptr(true),
	}

	if err = s.userReportRepo.UpdateMany(ctx, command.ReportedUserId, userReportEntity); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 3. Send mail deactivate user account
	if err = sendto.SendTemplateEmail(
		[]string{userUpdated.Email},
		consts.HOST_EMAIL,
		"deactivate_account.html",
		map[string]interface{}{"email": userUpdated.Email},
		"Yourvibes deactivated account",
	); err != nil {
		return response.NewCustomError(response.ErrSendEmailOTP)
	}

	return nil
}

func (s *sUserReport) DeleteUserReport(
	ctx context.Context,
	command *userCommand.DeleteUserReportCommand,
) (err error) {
	// 1. Check exists
	userReportFound, err := s.userReportRepo.GetById(ctx, command.UserId, command.ReportedUserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if userReportFound == nil {
		return response.NewDataNotFoundError("user report not found")
	}

	// 2. Check if report is already handled
	if userReportFound.Status {
		return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
	}

	// 3. Delete report
	if err = s.userReportRepo.DeleteOne(ctx, command.UserId, command.ReportedUserId); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sUserReport) ActivateUserAccount(
	ctx context.Context,
	command *userCommand.ActivateUserAccountCommand,
) (err error) {
	// 1. Check exists
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return response.NewDataNotFoundError("user not found")
	}

	// 2. Check if user is already activate
	if userFound.Status {
		return response.NewCustomError(response.ErrCodeUserIsAlreadyActivated)
	}

	// 3. Update reported user status
	reportedUserUpdateEntity := &userEntity.UserUpdate{
		Status: pointer.Ptr(true),
	}

	if err = reportedUserUpdateEntity.ValidateUserUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.userRepo.UpdateOne(ctx, command.UserId, reportedUserUpdateEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Update reportedUser posts status
	postUpdateEntity := &postEntity.PostUpdate{
		Status: pointer.Ptr(true),
	}
	if err = postUpdateEntity.ValidatePostUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	conditions := map[string]interface{}{
		"user_id": command.UserId,
	}

	if err = s.postRepo.UpdateMany(ctx, conditions, postUpdateEntity); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 5. Update reportedUser comment status
	conditions = map[string]interface{}{
		"user_id": command.UserId,
	}

	updateData := map[string]interface{}{
		"status": true,
	}

	if err = s.commentRepo.UpdateMany(ctx, conditions, updateData); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 6. Delete report
	if err = s.userReportRepo.DeleteByUserId(ctx, command.UserId); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Send email to user
	if err = sendto.SendTemplateEmail(
		[]string{userFound.Email},
		consts.HOST_EMAIL,
		"activate_account.html",
		map[string]interface{}{"email": userFound.Email},
		"Yourvibes activated account",
	); err != nil {
		return response.NewCustomError(response.ErrSendEmailOTP)
	}

	return nil
}

func (s *sUserReport) GetDetailUserReport(
	ctx context.Context,
	query *userQuery.GetOneUserReportQuery,
) (result *userQuery.UserReportQueryResult, err error) {
	// 1. Get user report detail
	userReportEntity, err := s.userReportRepo.GetById(ctx, query.UserId, query.ReportedUserId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if userReportEntity == nil {
		return nil, response.NewDataNotFoundError("user report not found")
	}

	// 2. Map to result
	return &userQuery.UserReportQueryResult{
		UserReport: mapper.NewUserReportResult(userReportEntity),
	}, nil
}

func (s *sUserReport) GetManyUserReport(
	ctx context.Context,
	query *userQuery.GetManyUserReportQuery,
) (result *userQuery.UserReportQueryListResult, err error) {
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

	return &userQuery.UserReportQueryListResult{
		UserReports:    userReportResults,
		PagingResponse: paging,
	}, nil
}
