package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	reportCommand "github.com/poin4003/yourVibes_GoApi/internal/application/report/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/mapper"
	reportQuery "github.com/poin4003/yourVibes_GoApi/internal/application/report/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sReportFactory struct {
	reportRepo repo.IReportRepository
}

func NewReportFactoryImplment(
	reportRepo repo.IReportRepository,
) *sReportFactory {
	return &sReportFactory{
		reportRepo: reportRepo,
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
		if err = s.reportRepo.HandleUserReport(ctx, command.ReportId, command.AdminId); err != nil {
			return err
		}
		return nil
	case consts.POST_REPORT:
		if err = s.reportRepo.HandlePostReport(ctx, command.ReportId, command.AdminId); err != nil {
			return err
		}
		return nil
	case consts.COMMENT_REPORT:
		if err = s.reportRepo.HandleCommentReport(ctx, command.ReportId, command.AdminId); err != nil {
			return err
		}
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
		if err = s.reportRepo.ActivateUser(ctx, command.ReportId); err != nil {
			return err
		}
		return nil
	case consts.POST_REPORT:
		if err = s.reportRepo.ActivatePost(ctx, command.ReportId); err != nil {
			return err
		}
		return nil
	case consts.COMMENT_REPORT:
		if err = s.reportRepo.ActivateComment(ctx, command.ReportId); err != nil {
			return err
		}
		return nil
	default:
		return response.NewValidateError("invalid report type")
	}
}
